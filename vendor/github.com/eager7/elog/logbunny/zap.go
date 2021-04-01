package logbunny

import (
	"errors"
	"math"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implement the logbunny logger
type zapLogger struct {
	lg                  *zap.Logger
	dynamicLevelHandler zap.AtomicLevel
	withCaller          bool
}

// this two func used for trans the level between zap & logbunny
func toZapLevel(lv LogLevel) zapcore.Level {
	switch lv {
	case DebugLevel:
		return zap.DebugLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case PanicLevel:
		return zap.PanicLevel
	case FatalLevel:
		return zap.FatalLevel
	default:
		return 0x7f
	}
}

func fromZapLevel(lv zapcore.Level) LogLevel {
	switch lv {
	case zap.DebugLevel:
		return DebugLevel
	case zap.InfoLevel:
		return InfoLevel
	case zap.WarnLevel:
		return WarnLevel
	case zap.ErrorLevel:
		return ErrorLevel
	case zap.PanicLevel:
		return PanicLevel
	case zap.FatalLevel:
		return FatalLevel
	default:
		return 0x7f
	}
}

// Level will also return a level handler for you, you can totally ignore it
func (log *zapLogger) Level() LogLevel {
	return fromZapLevel(log.dynamicLevelHandler.Level())
}

// SetLevel will also return a level handler for you, you can totally ignore it
func (log *zapLogger) SetLevel(lv LogLevel) {
	log.dynamicLevelHandler.SetLevel(toZapLevel(lv))
}

func (log *zapLogger) Log(lv LogLevel, msg string, fields ...*Field) {
	var msgFields []zapcore.Field
	msgFields = zapFields(fields...)

	switch lv {
	case DebugLevel:
		log.lg.Debug(msg, msgFields...)
	case InfoLevel:
		log.lg.Info(msg, msgFields...)
	case WarnLevel:
		log.lg.Warn(msg, msgFields...)
	case ErrorLevel:
		log.lg.Error(msg, msgFields...)
	case PanicLevel:
		log.lg.Panic(msg, msgFields...)
	case FatalLevel:
		log.lg.Fatal(msg, msgFields...)
	default:
		internalError("error in Log switch level", ErrTypeMissmatch, lv)
	}
}

func (log *zapLogger) Debug(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(DebugLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *zapLogger) Info(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(InfoLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *zapLogger) Warn(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(WarnLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *zapLogger) Error(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(ErrorLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *zapLogger) Panic(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(PanicLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *zapLogger) Fatal(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(FatalLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

// zapFields make zap field slice with given field
func zapFields(fields ...*Field) []zapcore.Field {
	msgfields := make([]zapcore.Field, 0, 3)
	for _, f := range fields {
		switch f.valType {
		case boolType:
			if f.vInt == 0 {
				msgfields = append(msgfields, zap.Bool(f.key, false))
			} else {
				msgfields = append(msgfields, zap.Bool(f.key, true))
			}
		case floatType:
			msgfields = append(msgfields, zap.Float64(f.key, math.Float64frombits(uint64(f.vInt))))
		case intType:
			msgfields = append(msgfields, zap.Int(f.key, int(f.vInt)))
		case int64Type:
			msgfields = append(msgfields, zap.Int64(f.key, f.vInt))
		case uintType:
			msgfields = append(msgfields, zap.Uint(f.key, uint(f.vInt)))
		case uint64Type:
			msgfields = append(msgfields, zap.Uint64(f.key, uint64(f.vInt)))
		case uintptrType:
			msgfields = append(msgfields, zap.Uintptr(f.key, uintptr(f.vInt)))
		case stringType:
			msgfields = append(msgfields, zap.String(f.key, f.vString))
		case timeType:
			msgfields = append(msgfields, zap.Time(f.key, f.vObj.(time.Time)))
		case durationType:
			msgfields = append(msgfields, zap.Duration(f.key, time.Duration(f.vInt)))
		case objectType:
			msgfields = append(msgfields, zap.Any(f.key, f.vObj))
		case stringerType:
			msgfields = append(msgfields, zap.String(f.key, f.vString))
		case errorType:
			msgfields = append(msgfields, zap.Error(errors.New(f.vString)))
		case skipType:
			continue
		default:
			internalError("unknown field type found", ErrTypeMissmatch, f)
		}
	}
	return msgfields
}
