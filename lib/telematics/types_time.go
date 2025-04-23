package telematics

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Unix timestamp for Apple reference time (2001)
const appleTimestampDiff int64 = 978307200

type Timestamp struct {
	time.Time
}

func (a *Timestamp) UnmarshalBSONValue(t bson.Type, raw []byte) (err error) {
	var val int64
	err = bson.UnmarshalValue(t, raw, &val)
	if err == nil {
		unix := val + appleTimestampDiff
		*a = Timestamp{time.Unix(unix, 0)}
	}
	return
}

type Duration struct {
	time.Duration
}

func (s *Duration) UnmarshalBSONValue(t bson.Type, raw []byte) (err error) {
	var val float64
	err = bson.UnmarshalValue(t, raw, &val)
	if err == nil {
		duration := val * float64(time.Second)
		*s = Duration{time.Duration(duration)}
	}
	return
}
