package db

type PayloadSize int

const (
	Small  PayloadSize = 1 << 8
	Medium PayloadSize = 1 << 9
	Large  PayloadSize = 1 << 10
	XLarge PayloadSize = 1 << 12
)

func GetPayloadBySizeKb(size PayloadSize) []byte {
	// TODO: generate payload of requested size
	return nil
}

func GetRandKeyUpTo(limit int) []byte {
	// TODO: generate a rand integer, from 1 up to limit and return []byte representation in 64B
	return nil
}
