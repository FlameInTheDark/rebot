package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

//CreateLogger creates production logger
func CreateLogger() *zap.Logger {
	w := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	return zap.New(core)
}

// CreateLoggerDebug creates debug logger
func CreateLoggerDebug() *zap.Logger {
	w := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.DebugLevel,
	)
	return zap.New(core)
}
