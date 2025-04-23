package telematics

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Fuzz(data []byte) int {
	var file File
	err := bson.Unmarshal(data, &file)

	if err == nil {
		return 1
	}

	return 0
}
