package overlay

import (
	"crypto/sha1"
	"encoding/binary"
)

func hash(ip string, capacity uint64) uint64 {
	// This has 160 bits.
	hashBytes := sha1.Sum(([]byte)(ip))
	relHashBytes := hashBytes[12:]

	relHashUint := binary.BigEndian.Uint64(relHashBytes)

	return relHashUint % (1 << capacity)
}

func isBetween(candidate uint64, start uint64, end uint64) bool {
	return ((start < end && candidate > start && candidate <= end) ||
		(start > end && (candidate > start || candidate <= end)))
}
