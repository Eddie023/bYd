package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/eddie023/byd/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("test info level", func(t *testing.T) {
		var out bytes.Buffer
		ctx := context.Background()
		log := logger.New(&out, "dummy-service")

		log.Info(ctx, "test info level", "foo", "bar")

		type logJSON struct {
			Level   string `json:"level"`
			File    string `json:"file"`
			Msg     string `json:"msg"`
			Service string `json:"service"`
		}
		var l logJSON
		err := json.Unmarshal([]byte(out.Bytes()), &l)
		require.NoError(t, err)

		assert.Contains(t, l.File, "logger_test.go:")
		assert.Equal(t, slog.LevelInfo.String(), l.Level)
		assert.Equal(t, l.Msg, "test info level")
		assert.Equal(t, l.Service, "dummy-service")
	})
}
