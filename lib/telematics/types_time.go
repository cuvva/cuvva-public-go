package telematics

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Unix timestamp for Apple reference time (2001)
const appleTimestampDiff int64 = 978307200

type Timestamp struct {
	time.Time
}

func (a *Timestamp) UnmarshalBSONValue(t bsontype.Type, raw []byte) (err error) {
	val, err := bsonrw.NewBSONValueReader(t, raw).ReadInt64()
	if err == nil {
		unix := val + appleTimestampDiff
		*a = Timestamp{time.Unix(unix, 0)}
	}
	return
}

type Duration struct {
	time.Duration
}

func (s *Duration) UnmarshalBSONValue(t bsontype.Type, raw []byte) (err error) {
	val, err := bsonrw.NewBSONValueReader(t, raw).ReadDouble()
	if err == nil {
		duration := val * float64(time.Second)
		*s = Duration{time.Duration(duration)}
	}
	return
}
