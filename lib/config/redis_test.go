package config

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedisOptions(t *testing.T) {
	expected := &redis.Options{
		Network:  "tcp",
		Addr:     "localhost:6379",
		Password: "password",
		DB:       1,
	}

	r := Redis{
		URI: "redis://:password@localhost/1",
	}

	opts, err := r.Options()

	assert.Nil(t, err)
	assert.Equal(t, expected, opts)
}
