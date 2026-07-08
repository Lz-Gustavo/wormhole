package db_test

import (
	"strconv"
	"testing"

	"github.com/Lz-Gustavo/wormhole/db"
	"github.com/stretchr/testify/assert"
)

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
			assert.Greater(t, len(key), 0)

			assert.GreaterOrEqual(t, n, int64(1))
			assert.LessOrEqual(t, n, tt.limit)
		})
	}
}
