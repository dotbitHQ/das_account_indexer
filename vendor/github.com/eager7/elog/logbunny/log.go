// Package logbunny is a go log framework warmed with serval popular logger.
// It is designed to take place of the slow old-fashion seelog.
// It's soooo powerful quick and flexible that everyone can't just believe it is called bunny ?!?
package logbunny

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Defined the common error from logbunny
var (
	ErrTypeMissmatch   = errors.New("error in convertion, type missmatch")
	ErrArgumentIllegal = errors.New("arguments illegal")
	ErrAssertion       = errors.New("error in assert the type")
	ErrConfigureError  = errors.New("error in apply the config, type miss match or illegal")
	ErrTeeLogger       = errors.New("error in tee the log, make sure the logger is legal & combining different logger is forbiden")
	ErrInternalError   = errors.New("error internal")
)

func internalError(msg string, err error, v interface{}) {
	fmt.Printf("Error internal ,%v ,%v , %v.\n", msg, err, v)
}

// Logger defined the interface which the logger need to implement
type Logger interface {
	// Debug and below functions is the logger family functions used to logout
	// the Level log
	Debug(string, ...*Field)
	Info(string, ...*Field)
	Warn(string, ...*Field)
	Error(string, ...*Field)
	Panic(string, ...*Field)
	Fatal(string, ...*Field)

	// Log will out put the log with given level. this will not make the pool recycle the field
	Log(LogLevel, string, ...*Field)

	// SetLevel() set a level into the logger Level() will return it
	SetLevel(LogLevel)
	Level() LogLevel
}

// New generate the logger implement the Logger with the given Options
func New(opts ...Option) (Logger, error) {
	cfg := &Config{
		Type:        ZapLogger,
		Level:       InfoLevel,
		Encoder:     JSONEncoder,
		WithCaller:  false,
		Out:         os.Stdout,
		TimePattern: "2006-01-02 15:04:05.999999999",
	}
	for _, v := range opts {
		v(cfg)
	}
	log, err := NewWithConfig(cfg)
	if err != nil {
		return nil, err
	}
	return log, nil
}

// NewWithConfig will generate with the given Config
func NewWithConfig(c *Config) (Logger, error) {
	// switch the type:
	// tzap: zap logger
	// tlogrus: logrus logger
	return newZapLogger(c)
}

// new Zap logger used to generate the zap logger implement logbunny logger
func newZapLogger(c *Config) (Logger, error) {
	var lv = zap.NewAtomicLevel()
	var encoder zapcore.Encoder
	var logger = &zapLogger{withCaller: c.WithCaller}
	var output zapcore.WriteSyncer

	if c.Out == nil {
		return nil, ErrConfigureError
	}

	output = zapcore.AddSync(c.Out)
	if !c.WithNoLock {
		output = zapcore.Lock(output)
	}

	switch c.Level {
	case DebugLevel:
		lv.SetLevel(zap.DebugLevel)
	case InfoLevel:
		lv.SetLevel(zap.InfoLevel)
	case WarnLevel:
		lv.SetLevel(zap.WarnLevel)
	case ErrorLevel:
		lv.SetLevel(zap.ErrorLevel)
	case PanicLevel:
		lv.SetLevel(zap.PanicLevel)
	case FatalLevel:
		lv.SetLevel(zap.FatalLevel)
	default:
		return nil, ErrConfigureError
	}
	logger.dynamicLevelHandler = lv

	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format(c.TimePattern))
	}
	encoderCfg := zapcore.EncoderConfig{
		NameKey:        "Name",
		StacktraceKey:  "Stack",
		MessageKey:     "Message",
		LevelKey:       "Level",
		TimeKey:        "TimeStamp",
		CallerKey:      "Caller",
		EncodeTime:     timeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var name string
	switch c.Encoder {
	case JSONEncoder:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
		name = "ZapJSONLogger"
	case ConsoleEncoder:
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		name = "ZapConsoleLogger"
	default:
		return nil, ErrConfigureError
	}

	log := zap.New(zapcore.NewCore(encoder, output, lv))
	//if c.WithCaller {
	//	log = zap.New(zapcore.NewCore(encoder, output, lv), zap.AddCaller())
	//}

	logger.lg = log
	log.Named(name)

	return logger, nil
}
