package cronzap

import (
	"go.uber.org/zap"
)

type Adapter struct {
	Log *zap.Logger
}

func (z *Adapter) Info(msg string, keysAndValues ...interface{}) {
	fields, ok := z.makeFields(keysAndValues)
	if ok {
		return
	}
	z.Log.Info(msg, fields...)
}

func (z *Adapter) makeFields(keysAndValues []interface{}) ([]zap.Field, bool) {
	if len(keysAndValues)%2 != 0 {
		z.Log.Warn("key / value pairs must be even.")
		return nil, false
	}

	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		if key, ok := keysAndValues[i].(string); ok {
			fields = append(fields, zap.Reflect(key, keysAndValues[i+1]))
		}

	}
	return fields, false
}

func (z *Adapter) Error(err error, msg string, keysAndValues ...interface{}) {
	fields, ok := z.makeFields(keysAndValues)
	if ok {
		return
	}

	z.Log.Error(msg, append(fields, zap.Error(err))...)
}
