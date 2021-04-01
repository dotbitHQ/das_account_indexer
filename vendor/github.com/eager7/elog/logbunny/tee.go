package logbunny

// Tee is used to combine some logger into one. Tee give the msg a copy to the sublogger to logout
func Tee(cfg *Config, logs ...Logger) Logger {
	switch len(logs) {
	case 0:
		return nil
	case 1:
		return logs[0]
	default:
		return &multiLogger{
			loggers:    logs,
			level:      cfg.Level,
			withCaller: cfg.WithCaller,
		}
	}
}

// multiLogger is used for Tee to combine Loggers into one
type multiLogger struct {
	loggers    []Logger
	level      LogLevel
	withCaller bool
}

func (log *multiLogger) SetLevel(lv LogLevel) {
	log.level = lv
	for _, lg := range log.loggers {
		lg.SetLevel(lv)
	}
}

func (log *multiLogger) Level() LogLevel {
	return log.level
}

// line below implement the logger
func (log *multiLogger) Log(lv LogLevel, msg string, fields ...*Field) {
	if log.level.enable(lv) {
		for _, lg := range log.loggers {
			lg.Log(lv, msg, fields...)
		}
	}
}
func (log *multiLogger) Debug(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(DebugLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}
func (log *multiLogger) Info(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(InfoLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}
func (log *multiLogger) Warn(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(WarnLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}
func (log *multiLogger) Error(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(ErrorLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}
func (log *multiLogger) Panic(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(PanicLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}
func (log *multiLogger) Fatal(msg string, fields ...*Field) {
	if log.withCaller {
		fields = append(fields, caller())
	}
	log.Log(FatalLevel, msg, fields...)
	for _, v := range fields {
		_fieldPool.Put(v)
	}
}
