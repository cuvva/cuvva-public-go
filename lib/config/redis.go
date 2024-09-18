package config

import (
	"time"

	"github.com/go-redis/redis"
)

// Redis configures a connection to a Redis database.
type Redis struct {
	URI          string        `json:"uri"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// Options returns a configured redis.Options structure.
func (r Redis) Options() (*redis.Options, error) {
	opts, err := redis.ParseURL(r.URI)
	if err != nil {
		return nil, err
	}

	opts.DialTimeout = r.DialTimeout
	opts.ReadTimeout = r.ReadTimeout
	opts.WriteTimeout = r.WriteTimeout

	return opts, nil
}

// Connect returns a connected redis.Client instance.
func (r Redis) Connect() (*redis.Client, error) {
	opts, err := r.Options()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping().Err(); err != nil {
		return client, err
	}

	return client, nil
}
