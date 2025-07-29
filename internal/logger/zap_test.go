package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInitializeLogger(t *testing.T) {
	// сбрасываем Log в zap.NewNop() для чистоты
	Log = zap.NewNop()

	err := Load()
	require.NoError(t, err, "logger initialization should not return error")

	// Проверим, что Log не остался zap.NewNop()
	nopLogger := zap.NewNop()
	assert.NotEqual(t, fmt.Sprintf("%p", nopLogger), fmt.Sprintf("%p", Log), "expected Log to be initialized, not zap.NewNop")

	// Дополнительно: пробуем что-то залогировать (не обязательно)
	Log.Info("test log message")
}
