package elog

import (
	"flag"
	"fmt"
	"github.com/eager7/elog/bunnystub"
	"github.com/eager7/elog/logbunny"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
)

var defaultLog logbunny.Logger //default logger
var loggerEncoder logbunny.EncoderType

// Zap & logrus level is not the same so we defined our level
const (
	NoticeLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var gLevel = NoticeLevel

const (
	colorRed = iota + 91
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
)

type loggerModule struct {
	name  string
	level int //debug-0 info-1 warn-2 err-3 panic-4 fatal-5
}

type Logger interface {
	Notice(a ...interface{})
	Debug(a ...interface{})
	Info(a ...interface{})
	Warn(a ...interface{})
	Error(a ...interface{})
	Fatal(a ...interface{})
	Panic(a ...interface{})
	GetLogger() logbunny.Logger
	SetLogLevel(level int) error
	GetLogLevel() int
	ErrStack()
}

func NewLogger(module string, level int) *loggerModule {
	return &loggerModule{name: module, level: level}
}

func init() {
	if err := Initialize(); err != nil {
		panic(err)
	}
}

func Initialize() error {
	var dir string
	if flag.Lookup("test.v") == nil {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		dir = strings.Replace(dir, "\\", "/", -1)
	} else {
		dir = "/tmp"
	}
	process := os.Args[0][strings.LastIndex(os.Args[0], `/`)+1:]
	fmt.Println("process name:", process)
	logOpt := ReadLoggerOpt(dir + "/" + process + ".yaml")
	gLevel = int(logOpt.debugLevel)
	dir += "/log/" + "elog_"

	logFilename := map[logbunny.LogLevel]string{
		logbunny.DebugLevel: dir + logOpt.debugLogFilename,
		logbunny.InfoLevel:  dir + logOpt.infoLogFilename,
		logbunny.WarnLevel:  dir + logOpt.warnLogFilename,
		logbunny.ErrorLevel: dir + logOpt.errorLogFilename,
		logbunny.FatalLevel: dir + logOpt.fatalLogFilename,
	}

	outputWriter := make(map[logbunny.LogLevel][]io.Writer)
	for level, file := range logFilename {
		var writer []io.Writer
		if logOpt.toTerminal {
			writer = append(writer, os.Stdout)
		}
		if logOpt.toFile {
			logFileWriter, err := newLogFile(file, logOpt.rollingTimePattern)
			if err != nil {
				return err
			}
			writer = append(writer, logFileWriter)
		}
		outputWriter[level] = writer
	}

	loggerEncoder = logOpt.loggerEncoder
	zapCfg := &logbunny.Config{
		Type:        logOpt.loggerType,
		Level:       logOpt.debugLevel,
		Encoder:     logOpt.loggerEncoder,
		WithCaller:  logOpt.withCaller,
		Out:         nil,
		WithNoLock:  false,
		TimePattern: logOpt.timePattern,
		Skip:        logOpt.skip,
	}
	l, err := logbunny.FilterLogger(zapCfg, outputWriter)
	if err != nil {
		return err
	}

	logbunny.SetCallerSkip(3)
	// log.Warp()
	defaultLog = l
	return nil
}

func newLogFile(logPath string, rollingTimePattern string) (io.Writer, error) {
	if file := stdOutput(logPath); file != nil {
		return file, nil
	} else {
		file, err := bunnystub.NewIOWriter(logPath, bunnystub.TimeRotate, bunnystub.WithTimePattern(rollingTimePattern))
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}
func stdOutput(logPath string) *os.File {
	if logPath == "stdout" {
		return os.Stdout
	}
	if logPath == "stderr" {
		return os.Stderr
	}
	return nil
}

func (l *loggerModule) Notice(a ...interface{}) {
	if l.level > NoticeLevel || gLevel > NoticeLevel {
		return
	}
	if loggerEncoder == 0 {
		defaultLog.Debug(fmt.Sprint(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorGreen) + "m" + "▶ " + "[" + l.name + "] " + fmt.Sprint(a...) + "\x1b[0m"
		defaultLog.Debug(msg)
	}
}

func (l *loggerModule) Debug(a ...interface{}) {
	if l.level > DebugLevel || gLevel > DebugLevel {
		return
	}
	if loggerEncoder == 0 {
		defaultLog.Debug(fmt.Sprint(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorBlue) + "m" + "▶ " + "[" + l.name + "] " + fmt.Sprint(a...) + "\x1b[0m"
		defaultLog.Debug(msg)
	}
}

func (l *loggerModule) Info(a ...interface{}) {
	if l.level > InfoLevel || gLevel > InfoLevel {
		return
	}
	if loggerEncoder == 0 {
		defaultLog.Info(fmt.Sprint(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorYellow) + "m" + "▶ " + "[" + l.name + "] " + fmt.Sprint(a...) + "\x1b[0m"
		defaultLog.Info(msg)
	}
}

func (l *loggerModule) Warn(a ...interface{}) {
	if l.level > WarnLevel || gLevel > WarnLevel {
		return
	}
	if loggerEncoder == 0 {
		defaultLog.Warn(fmt.Sprint(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorMagenta) + "m" + "▶ " + "[" + l.name + "] " + fmt.Sprint(a...) + "\x1b[0m"
		defaultLog.Warn(msg)
	}
}

func (l *loggerModule) Error(a ...interface{}) {
	if l.level > ErrorLevel || gLevel > ErrorLevel {
		return
	}
	if loggerEncoder == 0 {
		defaultLog.Error(fmt.Sprint(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorRed) + "m" + "▶ " + "[" + l.name + "] " + fmt.Sprint(a...) + "\x1b[0m"
		defaultLog.Error(msg)
	}
}

func (l *loggerModule) Fatal(a ...interface{}) {
	if l.level > FatalLevel || gLevel > FatalLevel {
		return
	}
	if loggerEncoder == 0 {
		defaultLog.Fatal(fmt.Sprint(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorYellow) + "m" + "▶ " + "[" + l.name + "] " + fmt.Sprint(a...) + "\x1b[0m"
		defaultLog.Fatal(msg)
	}
}

func (l *loggerModule) Panic(a ...interface{}) {
	if loggerEncoder == 0 {
		defaultLog.Panic(fmt.Sprint(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorYellow) + "m" + "▶ " + "[" + l.name + "] " + fmt.Sprint(a...) + "\x1b[0m"
		defaultLog.Panic(msg)
	}
	panic(fmt.Sprint(a...))
}

func (l *loggerModule) ErrStack() {
	l.Warn(string(debug.Stack()))
}

func (l *loggerModule) SetLogLevel(level int) error {
	l.level = level
	defaultLog.SetLevel(logbunny.LogLevel(level))
	return nil
}

func (l *loggerModule) LogLevel() int {
	return l.level
}

func (l *loggerModule) GlobalLevel() int {
	return gLevel
}

func (l *loggerModule) LogMode() bool {
	return l.level >= gLevel
}

func (l *loggerModule) Logger() logbunny.Logger {
	return defaultLog
}

func (l *loggerModule) WriteOptionToFile() {
	WriteLoggerOpt()
}

func (l *loggerModule) Noticef(format string, a ...interface{}) {
	l.Notice(fmt.Sprintf(format, a...))
}

func (l *loggerModule) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *loggerModule) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *loggerModule) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *loggerModule) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}
