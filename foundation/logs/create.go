package logs

import "go.uber.org/zap"

func CreateLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func CreateLoggerDebug() (*zap.Logger, error) {
	return zap.NewDevelopment()
}