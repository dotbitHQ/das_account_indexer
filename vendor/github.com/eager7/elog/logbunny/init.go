package logbunny

var (
	// Log can used to print the log directly
	Log Logger
)

func init() {
	var err error
	Log, err = New()
	if err != nil {
		internalError("error in init the logger", err, nil)
	}
}

// Debug can print out the log directly with the Log
func Debug(msg string, fields ...*Field) {
	_callerSkip = 3
	Log.Debug(msg, fields...)
}

// Info can print out the log directly with the Log
func Info(msg string, fields ...*Field) {
	_callerSkip = 3
	Log.Info(msg, fields...)
}

// Warn can print out the log directly with the Log
func Warn(msg string, fields ...*Field) {
	_callerSkip = 3
	Log.Warn(msg, fields...)
}

// Error can print out the log directly with the Log
func Error(msg string, fields ...*Field) {
	_callerSkip = 3
	Log.Error(msg, fields...)
}

// Panic can print out the log directly with the Log
func Panic(msg string, fields ...*Field) {
	_callerSkip = 3
	Log.Panic(msg, fields...)
}

// Fatal can print out the log directly with the Log
func Fatal(msg string, fields ...*Field) {
	_callerSkip = 3
	Log.Fatal(msg, fields...)
}
