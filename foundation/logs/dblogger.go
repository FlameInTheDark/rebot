package logs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sqldblogger "github.com/simukti/sqldb-logger"
	"go.uber.org/zap"
)

const (
	LongDataFields = `query args`
	MaxDataLen     = 200
)

var _ sqldblogger.Logger = (*DBLogger)(nil)

type DBLogger struct {
	log *zap.Logger
}

func NewDBLogger(log *zap.Logger) *DBLogger {
	return &DBLogger{log: log.With(zap.String("module", "database"))}
}

func (dbl *DBLogger) Log(_ context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	dbl.trimLongData(data)
	switch level {
	case sqldblogger.LevelError:
		dbl.log.Error(msg, zap.Reflect("data", data))
	case sqldblogger.LevelInfo:
		dbl.log.Info(msg, zap.Reflect("data", data))
	case sqldblogger.LevelDebug, sqldblogger.LevelTrace:
		fallthrough
	default:
		dbl.log.Debug(msg, zap.Reflect("data", data))
	}
}

// trimLongData trims long data values with keys in "LongDataFields"
func (z *DBLogger) trimLongData(data map[string]interface{}) {
	for k, v := range data {
		if !strings.Contains(LongDataFields, k) {
			continue
		}

		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		if length := len(str); length > MaxDataLen {
			data[k] = str[:MaxDataLen] + "... (" + strconv.Itoa(length) + " symbols)"
		}
	}
}
