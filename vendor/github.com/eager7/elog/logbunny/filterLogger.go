package logbunny

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	colorRed = iota + 91
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
)

// filterLogger implement the Logger and contains sublogger which
// ONLY logout the level filter log into given io.Writer
type filterLogger struct {
	loggers    map[LogLevel][]Logger
	level      LogLevel
	withCaller bool
}

// FilterLogger will generate the log with given outs & config
// Split the different log into sub logger with io.Writer
func FilterLogger(c *Config, outs map[LogLevel][]io.Writer) (Logger, error) {
	subLoggers := make(map[LogLevel][]Logger)
	for k, v := range outs {
		for _, out := range v {
			cfg := *c
			cfg.Out = out
			log, err := newZapFilterLogger(&cfg, k)
			if err != nil {
				return nil, err
			}
			subLoggers[k] = append(subLoggers[k], log)
		}
	}
	return &filterLogger{loggers: subLoggers, level: c.Level, withCaller: c.WithCaller}, nil
}

// new zap filter logger will give out an zap logger with given config , this logger level cannot be changed
func newZapFilterLogger(c *Config, lvl LogLevel) (Logger, error) {
	var encoder zapcore.Encoder
	var logger = &zapLogger{}
	var output zapcore.WriteSyncer
	var lvebf zap.LevelEnablerFunc
	if c.Out == nil {
		return nil, ErrConfigureError
	}
	output = zapcore.AddSync(c.Out)
	zlv := toZapLevel(lvl)
	lvebf = zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l == zlv })

	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format(c.TimePattern))
	}
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "Msg",
		LevelKey:       "Lvl",
		TimeKey:        "TS",
		CallerKey:      "Caller",
		EncodeTime:     timeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	switch c.Encoder {
	case JSONEncoder:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	case ConsoleEncoder:
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	default:
		return nil, ErrConfigureError
	}

	var log *zap.Logger
	if c.WithCaller {
		log = zap.New(zapcore.NewCore(encoder, output, lvebf), zap.AddCaller(), zap.AddCallerSkip(c.Skip))
	} else {
		log = zap.New(zapcore.NewCore(encoder, output, lvebf), zap.AddCaller(), zap.AddCallerSkip(c.Skip))
	}
	log.Named("EcoballLogger")

	logger.lg = log

	return logger, nil
}

// SetLevel set the filter logger level with given level
// SetLevel will not and should not apply the level to the subloggers
func (log *filterLogger) SetLevel(lv LogLevel) {
	log.level = lv
}

// Level return the filter logger level , this level is a global level filter
// which control the output weather call the sub logger to log the message
func (log *filterLogger) Level() LogLevel {
	return log.level
}

// Line below implement the Logger, only different is here got a enabler to check the level
func (log *filterLogger) Log(lv LogLevel, msg string, fields ...*Field) {
	if log.level.enable(lv) {
		for _, lg := range log.loggers[lv] {
			lg.Log(lv, msg, fields...)
		}
	}
}

func (log *filterLogger) Debug(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	//log.Log(DebugLevel, "\x1b[" + strconv.Itoa(colorBlue) + "m"  + msg + "\x1b[0m ", fields...)
	log.Log(DebugLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *filterLogger) Info(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	//log.Log(InfoLevel, "\x1b[" + strconv.Itoa(colorYellow) + "m"  + msg + "\x1b[0m ", fields...)
	log.Log(InfoLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *filterLogger) Warn(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	//log.Log(WarnLevel, "\x1b[" + strconv.Itoa(colorMagenta) + "m"  + msg + "\x1b[0m ", fields...)
	log.Log(WarnLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *filterLogger) Error(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	//log.Log(ErrorLevel, "\x1b[" + strconv.Itoa(colorRed) + "m"  + msg + "\x1b[0m ", fields...)
	log.Log(ErrorLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *filterLogger) Panic(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	//log.Log(PanicLevel, "\x1b[" + strconv.Itoa(colorRed) + "m"  + msg + "\x1b[0m ", fields...)
	log.Log(PanicLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}

func (log *filterLogger) Fatal(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	//log.Log(FatalLevel, "\x1b[" + strconv.Itoa(colorRed) + "m"  + msg + "\x1b[0m ", fields...)
	log.Log(FatalLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}
