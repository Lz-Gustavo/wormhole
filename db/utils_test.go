package db_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/stretchr/testify/assert"
)

func Test_GetPayloadBySizeKb(t *testing.T) {
	tests := []struct {
		name            string
		size            db.PayloadSize
		expectedPayload string
	}{
		{
			name:            "successfully generate valid payload",
			size:            db.Small,
			expectedPayload: strings.Repeat("0", int(db.Small)),
		},
		{
			name:            "failed generating failed and return empty string",
			size:            db.PayloadSize(10),
			expectedPayload: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := db.GetPayloadBySizeKb(tt.size)

			assert.Equal(t, tt.expectedPayload, v)
			if tt.expectedPayload != "" {
				assert.Len(t, v, int(tt.size))
			}
		})
	}
}

func Test_GetRandKeyUpTo(t *testing.T) {
	tests := []struct {
		name        string
		limit       int64
		shouldPanic bool
	}{
		{
			name:        "successfully generate rand key within limit",
			limit:       10,
			shouldPanic: false,
		},
		{
			name:        "successfully generate rand key within greater limit",
			limit:       2189021217621,
			shouldPanic: false,
		},
		{
			name:        "panic on non-positive limit",
			limit:       -1,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() { db.GetRandKeyUpTo(tt.limit) })
				return
			}
			key := db.GetRandKeyUpTo(tt.limit)

			n, err := strconv.ParseInt(key, 10, 64)
			assert.NoError(t, err)
			assert.Len(t, key, db.KeySizeBytes)

			assert.GreaterOrEqual(t, n, int64(1))
			assert.LessOrEqual(t, n, tt.limit)
		})
	}
}
