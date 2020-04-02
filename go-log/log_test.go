package log

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestGetLastErrors(t *testing.T) {
	var config = LogConfig{
		path:   "./",
		prefix: "demo",
		maxAge: 7,
	}

	initConfig(config)

	Logger.Info("Finished")
	Logger.Warn("Finished")
	Logger.Error("Game Over")
	SugarLogger.Info("Finished")

    SetDebugMod() //debug模式
    SugarLogger.Debug("Finished")

	SugarLogger.With(
		"hello", "world",
		"failure", errors.New("oh no"),
		"count", 42,
		"user", "alice",
	)
	SugarLogger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", "http://baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	SugarLogger.Infow("info w", map[string]interface{}{"test" : "test"})

	Logger.Error("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", "http://baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
