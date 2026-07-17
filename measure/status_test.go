package measure_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Lz-Gustavo/wormhole/measure"
	"github.com/stretchr/testify/assert"
)

func Test_StatusMsr(t *testing.T) {
	testExecTime := 2 * time.Second

	tmpDir := t.TempDir()
	fn := filepath.Join(tmpDir, "test-status.out")

	sm, err := measure.NewStatusMsr(fn)
	assert.NoError(t, err)
	defer sm.Close()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go sm.Run(ctx)

	countSuccess := 3
	countTimeout := 2
	countFail := 1

	for range countSuccess {
		sm.CountStatusFromErr(nil)
	}
	for range countTimeout {
		sm.CountStatusFromErr(context.DeadlineExceeded)
	}
	for range countFail {
		sm.CountStatusFromErr(errors.New("err"))
	}

	time.Sleep(testExecTime)
	cancel()

	assert.NoError(t, sm.Flush())
	assert.NoError(t, sm.Close())

	data, err := os.ReadFile(fn)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	totalSuccess := 0
	totalTimeout := 0
	totalFail := 0

	for _, line := range lines {
		if line == "" {
			continue
		}

		var gotSuccess, gotTimeout, gotFail int
		_, err := fmt.Sscanf(line, "200:%d-408:%d-500:%d\n", &gotSuccess, &gotTimeout, &gotFail)
		assert.NoError(t, err)

		totalSuccess += gotSuccess
		totalTimeout += gotTimeout
		totalFail += gotFail
	}

	assert.Equal(t, countSuccess, totalSuccess)
	assert.Equal(t, countTimeout, totalTimeout)
	assert.Equal(t, countFail, totalFail)
}
