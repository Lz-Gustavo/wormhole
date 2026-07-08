package db

type PayloadSize int

const (
	keySizeBytes = 64

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

func GetPayloadBySizeKb(size PayloadSize) string {
	// TODO: generate payload of requested size (already validated)
	return ""
}

func GetRandKeyUpTo(limit int64) string {
	// TODO: generate a rand integer, from 1 up to limit and return string representation
	// from a []byte of 64B
	return ""
}
