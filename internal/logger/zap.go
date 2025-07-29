package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log — глобальный логгер, используемый во всём приложении.
//
// Инициализируется функцией Load(). По умолчанию — zap.NewNop() (без вывода).
var Log *zap.Logger = zap.NewNop()

// Load инициализирует глобальный логгер Log с использованием конфигурации `zap.NewProductionConfig()`.
//
// Устанавливает уровень логирования DebugLevel.
// Возвращает ошибку в случае неудачи при построении логгера.
func Load() error {
	cfg := zap.NewProductionConfig()

	cfg.Level.SetLevel(zapcore.DebugLevel)

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
