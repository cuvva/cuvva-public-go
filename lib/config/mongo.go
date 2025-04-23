package config

import (
	"context"
	"errors"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/db/mongodb"
	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

// MongoDB configures a connection to a Mongo database.
type MongoDB struct {
	URI             string         `json:"uri"`
	ConnectTimeout  time.Duration  `json:"connect_timeout"`
	MaxConnIdleTime *time.Duration `json:"max_conn_idle_time"`
	MaxConnecting   *uint64        `json:"max_connecting"`
	MaxPoolSize     *uint64        `json:"max_pool_size"`
	MinPoolSize     *uint64        `json:"min_pool_size"`
}

// Options returns the MongoDB client options and database name.
func (m MongoDB) Options() (opts *options.ClientOptions, dbName string, err error) {
	opts = options.Client().ApplyURI(m.URI)
	opts.MaxConnIdleTime = m.MaxConnIdleTime
	opts.MaxConnecting = m.MaxConnecting
	opts.MaxPoolSize = m.MaxPoolSize
	opts.MinPoolSize = m.MinPoolSize
	if m.ConnectTimeout == 0 {
		m.ConnectTimeout = 10 * time.Second
	}
	opts.ConnectTimeout = ptr.Ptr(m.ConnectTimeout)

	err = opts.Validate()
	if err != nil {
		return
	}

	// all Go services use majority reads/writes, and this is unlikely to change
	// if it does change, switch to accepting as an argument
	opts.SetReadConcern(readconcern.Majority())
	opts.SetWriteConcern(&writeconcern.WriteConcern{
		W:       writeconcern.WCMajority,
		Journal: ptr.Bool(true),
	})

	cs, err := connstring.Parse(m.URI)
	if err != nil {
		return
	}

	dbName = cs.Database
	if dbName == "" {
		err = errors.New("missing mongo database name")
	}

	return
}

// Connect returns a connected mongo.Database instance.
func (m MongoDB) Connect() (*mongodb.Database, error) {
	opts, dbName, err := m.Options()
	if err != nil {
		return nil, err
	}

	return mongodb.Connect(opts, dbName)
}
