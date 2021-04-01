package elog

import (
	"fmt"
	"github.com/eager7/elog/logbunny"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

type loggerOpt struct {
	debugLevel         logbunny.LogLevel
	loggerType         logbunny.LogType
	withCaller         bool
	toTerminal         bool
	toFile             bool
	loggerEncoder      logbunny.EncoderType
	timePattern        string
	debugLogFilename   string
	infoLogFilename    string
	warnLogFilename    string
	errorLogFilename   string
	fatalLogFilename   string
	rollingTimePattern string
	skip               int
	logger             logbunny.Logger
}

func ReadLoggerOpt(file string) *loggerOpt {
	viper.SetConfigFile(file)
	initDefaultConfig()
	_ = viper.ReadInConfig()
	go watchConfig(file)

	return &loggerOpt{
		debugLevel:         logbunny.LogLevel(viper.GetInt("log.debug_level")),
		loggerType:         logbunny.LogType(viper.GetInt("log.loggerType")),
		withCaller:         viper.GetBool("log.with_caller"),
		toTerminal:         viper.GetBool("log.to_terminal"),
		toFile:             viper.GetBool("log.to_file"),
		loggerEncoder:      logbunny.EncoderType(viper.GetInt("log.logger_encoder")),
		timePattern:        viper.GetString("log.time_pattern"),
		debugLogFilename:   viper.GetString("log.debug_log_filename"),
		infoLogFilename:    viper.GetString("log.info_log_filename"),
		warnLogFilename:    viper.GetString("log.warn_log_filename"),
		errorLogFilename:   viper.GetString("log.error_log_filename"),
		fatalLogFilename:   viper.GetString("log.fatal_log_filename"),
		rollingTimePattern: viper.GetString("log.rolling_time_pattern"),
		skip:               viper.GetInt("log.skip"),
		logger:             nil,
	}
}

func WriteLoggerOpt() {
	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		if _, err := os.Create(viper.ConfigFileUsed()); err == nil {
			if err := viper.WriteConfig(); err != nil {
				fmt.Println("write err:", err)
			}
		} else {
			fmt.Println("create err:", err)
		}
	}
}

func initDefaultConfig() {
	viper.SetDefault("log.debug_level", 0) //0: debug 1: info 2: warn 3: error 4: panic 5: fatal
	viper.SetDefault("log.loggerType", 0)  //0: zap 1: logrus
	viper.SetDefault("log.with_caller", false)
	viper.SetDefault("log.logger_encoder", 1) //0: json 1: console
	viper.SetDefault("log.time_pattern", "2006-01-02 15:04:05.00000")
	viper.SetDefault("log.http_port", ":50015")                 // RESTFul API to change logout level dynamically
	viper.SetDefault("log.debug_log_filename", "debug.log")     //or 'stdout' / 'stderr'
	viper.SetDefault("log.info_log_filename", "info.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.warn_log_filename", "warn.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.error_log_filename", "err.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.fatal_log_filename", "fatal.log")     //or 'stdout' / 'stderr'
	viper.SetDefault("log.rolling_time_pattern", "0 0 0 * * *") //rolling the log everyday at 00:00:00
	viper.SetDefault("log.skip", 4)                             //call depth, zap log is 3, logger is 4
	viper.SetDefault("log.to_terminal", true)                   //out put to terminal
	viper.SetDefault("log.to_file", false)                      //out put to file
}

func watchConfig(file string) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		viper.SetConfigFile(file)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("retry read config file error:", err)
		} else {
			level := int(viper.GetInt("log.debug_level"))
			fmt.Println("new config:", level)
			if level >= 0 && level <= FatalLevel && level != gLevel {
				gLevel = level
			}
		}
	})
}
