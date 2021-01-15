package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ErrorCodeDuplicateKey = 11000
)

func HasErrorCode(err error, code int) bool {
	we, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}

	for _, err := range we.WriteErrors {
		if err.Code == code {
			return true
		}
	}

	return false
}
