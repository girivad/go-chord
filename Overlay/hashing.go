package overlay

import (
	"crypto/sha1"
	"encoding/binary"
)

func hash(ip string, capacity int64) int64 {
	// This has 160 bits.
	hashBytes := sha1.Sum(([]byte)(ip))
	relHashBytes := hashBytes[13:]

	relHashUint := binary.BigEndian.Uint64(relHashBytes)

	relHashInt := int64(relHashUint)

	return relHashInt % (1 << capacity)
}

func isBetween(candidate int64, start int64, end int64) bool {
	return ((start < end && candidate > start && candidate <= end) ||
		(start > end && (candidate > start || candidate <= end)))
}
