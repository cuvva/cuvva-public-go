package telematics

import (
	"encoding/binary"
	"errors"
	"math"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// sampling data is an efficient binary-packed format
// each []byte data field contains a set of samples (4 bytes per sample)
// each sample contains 3 values of 10 bits each - last 2 bits unused
// the 10-bit values are normalized to a unit-specific signed range

type GForce float64
type Radians float64

type Acceleration struct {
	X, Y, Z GForce
}

type Attitude struct {
	Roll, Pitch, Yaw Radians
}

const (
	maxAcceleration = GForce(10)
	maxAttitude     = Radians(math.Pi)

	// these apply to both acceleration and attitude
	valueNorm          = float64(500)
	valueBits          = 10
	valueNormMask      = ((1 << valueBits) - 1)    // 1023
	valueNormSignedMax = ((valueNormMask + 1) / 2) // 512
)

const sampleSize = 4 // 4 bytes in a uint32

type AccelerationData []Acceleration
type AttitudeData []Attitude

func (a *AccelerationData) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	count, data, err := prepSamplingData(t, raw)
	if err != nil {
		return err
	}

	out := make(AccelerationData, count)

	for i := 0; i < count; i++ {
		x, y, z := readSample(data, i, float64(maxAcceleration))

		out[i] = Acceleration{
			GForce(x),
			GForce(y),
			GForce(z),
		}
	}

	*a = out

	return nil
}

func (a *AttitudeData) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	count, data, err := prepSamplingData(t, raw)
	if err != nil {
		return err
	}

	out := make(AttitudeData, count)

	for i := 0; i < count; i++ {
		roll, pitch, yaw := readSample(data, i, float64(maxAttitude))

		out[i] = Attitude{
			Radians(roll),
			Radians(pitch),
			Radians(yaw),
		}
	}

	*a = out

	return nil
}

func prepSamplingData(t bsontype.Type, raw []byte) (count int, data []byte, err error) {
	data, _, err = bsonrw.NewBSONValueReader(t, raw).ReadBinary()
	if err != nil {
		return
	}

	if len(data)%sampleSize != 0 {
		err = errors.New("invalid []byte length")
		return
	}

	count = len(data) / sampleSize
	return
}

func readSample(data []byte, idx int, max float64) (a, b, c float64) {
	offset := idx * sampleSize
	end := offset + sampleSize
	packed := binary.BigEndian.Uint32(data[offset:end])

	return unpackSample(packed, max)
}

func unpackSample(packed uint32, max float64) (a, b, c float64) {
	a = denormalizeValue(packed>>(0*valueBits), max)
	b = denormalizeValue(packed>>(1*valueBits), max)
	c = denormalizeValue(packed>>(2*valueBits), max)

	return
}

func denormalizeValue(input uint32, max float64) float64 {
	value := int32(input) & valueNormMask

	if value >= valueNormSignedMax {
		value -= valueNormSignedMax * 2
	}

	return float64(value) * max / valueNorm
}
