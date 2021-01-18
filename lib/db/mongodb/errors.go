package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ErrorCodeDuplicateKey = 11000
)

func HasErrorCode(err error, code int) bool {
	switch v := err.(type) {
	case mongo.BulkWriteException:
		fallthrough
	case mongo.WriteException:
		for _, err := range v.WriteErrors {
			if err.Code == code {
				return true
			}
		}
	}

	return false
}
