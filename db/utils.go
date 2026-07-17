package db

import (
	"math/rand/v2"
	"strconv"
	"strings"
)

type PayloadSize int

const (
	KeySizeBytes = 64

	Small  PayloadSize = 1 << 8
	Medium PayloadSize = 1 << 9
	Large  PayloadSize = 1 << 10
	XLarge PayloadSize = 1 << 12
)

var validPayloadSizes = map[PayloadSize]struct{}{
	Small:  {},
	Medium: {},
	Large:  {},
	XLarge: {},
}

func (ps PayloadSize) IsValid() bool {
	_, exists := validPayloadSizes[ps]
	return exists
}

// GetPayloadBySizeKb returns a zero-filled payload of the requested size, or an empty string for invalid sizes.
func GetPayloadBySizeKb(size PayloadSize) string {
	if !size.IsValid() {
		return ""
	}
	return strings.Repeat("0", int(size))
}

// GetRandKeyUpTo returns a left-padded decimal key of length KeySizeBytes in the range [1, limit].
func GetRandKeyUpTo(limit int64) string {
	key := rand.Int64N(limit) + 1
	skey := strconv.FormatInt(key, 10)
	paddingLen := KeySizeBytes - len(skey)

	return strings.Repeat("0", paddingLen) + skey
}
