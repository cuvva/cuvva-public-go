package limiter

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// Sliding implements a Redis-backed Limiter for a sliding window.
//
// The sliding window is implemented using Redis sorted sets, where each
// event is a new entry in the sorted set with its timestamp as a score. Events
// older than the window are evicted and the set eventually expires.
type Sliding struct {
	// WindowDuration defines the width of the sliding window
	// where events are counted against the maximum.
	WindowDuration time.Duration

	// WindowMaximum is the maximum number of events that can
	// happen within the configured window.
	WindowMaximum int64

	// Redis is the storage backend for per-token rate limit counts.
	Redis redis.Cmdable

	// RedisPrefix will prefix all keys used by the Sliding limiter.
	RedisPrefix string
}

// Allow reports whether an even with the given token can happen
// within the configured maximum rate. Returned is the current
// event count and whether the event can happen.
//
// If redis is unavailable, Allow will allow all tokens temporarily.
func (s *Sliding) Allow(token string) (count int64, ok bool) {
	if s.RedisPrefix != "" {
		token = s.RedisPrefix + "/" + token
	}

	now := time.Now()
	prior := now.Add(-(s.WindowDuration))

	pipe := s.Redis.TxPipeline()

	// trim elements older than window from sorted set
	// get cardinality of elements still in set
	// add current second as element in set (if not already added)
	// reset expiry of sorted set to window
	pipe.ZRemRangeByScore(token, "0", strconv.FormatInt(prior.Unix(), 10))
	pipe.ZCard(token)
	pipe.ZAddNX(token, redis.Z{Score: float64(now.Unix()), Member: now.Unix()})
	pipe.Expire(token, s.WindowDuration)

	// TODO(jc): replace the above redis calls with lua, can then invoke SHA1
	// hash of script instead of sending command each time.

	result, err := pipe.Exec()
	if err != nil || len(result) != 4 {
		// allow all events when cannot connect to redis
		count = s.WindowMaximum
		ok = true
		return
	}

	if card, isOK := result[1].(*redis.IntCmd); isOK {
		count = card.Val()
		if count < s.WindowMaximum {
			ok = true
		}
	}

	return
}
