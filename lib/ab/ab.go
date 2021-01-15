package ab

import (
	"crypto/sha1"
)

const BucketCount = 65536

func Check(testID, subjectID string, maxBucketIndex uint16) bool {
	hash := sha1.Sum([]byte(testID + subjectID))
	bucket := uint16(hash[6])<<8 + uint16(hash[7]) // equivalent to uint64 % BucketCount

	return bucket <= maxBucketIndex
}
