package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

func NewRandFromString(str string) *rand.Rand {
	hash := sha256.New()
	hash.Write([]byte(str))

	seed := binary.BigEndian.Uint64(hash.Sum(nil))

	randSrc := rand.NewSource(int64(seed))

	return rand.New(randSrc)
}

func PercentageFromString(str string) float64 {
	hash := sha256.New()
	hash.Write([]byte(str))

	n := binary.BigEndian.Uint64(hash.Sum(nil))
	n = n % 100_000

	return float64(n) / 1000.0
}
