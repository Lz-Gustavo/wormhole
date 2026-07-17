package measure_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Lz-Gustavo/wormhole/measure"
	"github.com/stretchr/testify/assert"
)

func Test_LatencyMsr(t *testing.T) {
	tmpDir := t.TempDir()
	fn := filepath.Join(tmpDir, "test-latency.out")

	lm, err := measure.NewLatencyMsr(fn)
	assert.NoError(t, err)
	defer lm.Close()

	latencies := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
	}
	for _, lat := range latencies {
		err := lm.Record(lat)
		assert.NoError(t, err)
	}
	assert.NoError(t, lm.Flush())
	assert.NoError(t, lm.Close())

	data, err := os.ReadFile(fn)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	assert.Equal(t, len(latencies), len(lines))

	for i, line := range lines {
		lat, err := strconv.Atoi(line)
		assert.NoError(t, err)

		expected := int(latencies[i].Nanoseconds())
		assert.Equal(t, expected, lat)
	}
}
