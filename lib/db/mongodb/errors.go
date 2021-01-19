package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ErrorCodeDuplicateKey = 11000
)

func HasErrorCode(err error, code int) bool {
	var writeErrors []mongo.WriteError

	switch v := err.(type) {
	case mongo.WriteException:
		writeErrors = append(writeErrors, v.WriteErrors...)
	case mongo.BulkWriteException:
		for _, err := range v.WriteErrors {
			writeErrors = append(writeErrors, err.WriteError)
		}
	}

	for _, err := range writeErrors {
		if err.Code == code {
			return true
		}
	}

	return false
}
