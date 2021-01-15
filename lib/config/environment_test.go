package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type env map[string]string

func (e env) Get(k string) string {
	return e[k]
}

func TestFromEnvironment(t *testing.T) {
	env := env{ConfigEnvironmentVariable: `{"foo": "bar"}`}

	v := struct {
		Foo string `json:"foo"`
	}{}

	err := FromEnvironment(env.Get, &v)
	if assert.NoError(t, err) {
		assert.Equal(t, "bar", v.Foo)
	}
}

func ExampleFromEnvironment() {
	var config struct {
		CacheRedis Redis `json:"cache_redis"`
	}

	err := FromEnvironment(os.Getenv, &config)
	if err != nil {
		panic(err)
	}
}

func TestEnvironmentName(t *testing.T) {
	env := env{ConfigEnvironmentVariable: `{"env": "prod"}`}

	assert.Equal(t, "prod", EnvironmentName(env.Get))
}

func TestEnvironmentNameDev(t *testing.T) {
	env := env{}

	assert.Equal(t, "local", EnvironmentName(env.Get))
}
