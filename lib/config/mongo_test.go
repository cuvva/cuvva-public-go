package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

func TestMongoDBOptions(t *testing.T) {
	m := &MongoDB{
		URI: "mongodb://foo:bar@127.0.0.1/demo?authSource=admin",
	}

	opts, dbName, err := m.Options()

	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, dbName, "demo")
	assert.Equal(t, opts.Hosts, []string{"127.0.0.1"})
	assert.Equal(t, opts.WriteConcern, writeconcern.New(writeconcern.WMajority(), writeconcern.J(true)))

	assert.Equal(t, opts.Auth, &options.Credential{
		AuthSource:  "admin",
		Username:    "foo",
		Password:    "bar",
		PasswordSet: true,
	})
}
