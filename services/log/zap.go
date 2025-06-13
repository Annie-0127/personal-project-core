package log

import (
	"go.uber.org/zap"
	"os"
)

func NewZapLogger() *zap.Logger {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	logger := zap.Must(zap.NewDevelopment())
	if env == "production" {
		logger = zap.Must(zap.NewProduction())
	}

	defer logger.Sync()

	return logger
}
