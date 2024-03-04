package overlay

import (
	"crypto/sha1"
	"encoding/binary"
)

func hash(ip string, capacity int64) int64 {
	// Placeholder Hash
	// This has 160 bits.
	hashBytes := sha1.Sum(([]byte)(ip))
	relHashBytes := hashBytes[19-capacity/8:]

	return int64(binary.BigEndian.Uint64(relHashBytes)) % (1 << capacity)
}

func isBetween(candidate int64, start int64, end int64) bool {
	return ((start < end && candidate > start && candidate < end) ||
		(start > end && (candidate > start || candidate <= end)))
}
