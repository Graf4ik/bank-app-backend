package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	envLocal string = "local"
	envDev   string = "dev"
	envProd  string = "prod"
)

var Log *zap.Logger

func InitLogger(env string) *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	switch env {
	case envLocal:
		Log = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		))
	case envDev:
		Log, _ = zap.NewDevelopment()
	case envProd:
		Log, _ = zap.NewProduction()
	}

	return Log
}
