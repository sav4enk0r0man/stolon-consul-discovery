package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Logger struct {
	Fatal *log.Logger
	Error *log.Logger
	Info  *log.Logger
	Debug *log.Logger
	Trace *log.Logger
	Level LogLevel
}

type LogLevel struct {
	level string
}

var (
	FATAL   = LogLevel{"FATAL"}
	ERROR   = LogLevel{"ERROR"}
	INFO    = LogLevel{"INFO"}
	DEBUG   = LogLevel{"DEBUG"}
	TRACE   = LogLevel{"TRACE"}
	UNKNOWN = LogLevel{""}
)

var DefaultLog *Logger

func (l LogLevel) String() string {
	return l.level
}

func FromString(str string) (LogLevel, error) {
	switch strings.ToUpper(str) {
	case FATAL.level:
		return FATAL, nil
	case ERROR.level:
		return ERROR, nil
	case INFO.level:
		return INFO, nil
	case DEBUG.level:
		return DEBUG, nil
	case TRACE.level:
		return TRACE, nil
	}
	return UNKNOWN, fmt.Errorf("unknown log level: %s", str)
}

func (l LogLevel) Index() int {
	levels := []LogLevel{0: FATAL, 1: ERROR, 2: INFO, 3: DEBUG, 4: TRACE}
	for i, lvl := range levels {
		if lvl == l {
			return i
		}
	}
	return -1
}

func Wrapper(err error, msg string) error {
	if err != nil {
		pc, filename, line, _ := runtime.Caller(1)
		return fmt.Errorf("in %s[%s:%d]: \t%v",
			filepath.Base(runtime.FuncForPC(pc).Name()), filepath.Base(filename), line, msg)
	}
	return err
}

func NewLogger(conf map[string]string) *Logger {
	devNull := ioutil.Discard
	logFormat := conf["logformat"]

	se := os.Stderr
	so := os.Stdout

	if conf["logfile"] != "" {
		logfile, err := os.OpenFile(conf["logfile"], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Panic(err)
		}
		se = logfile
		so = logfile
	}

	if conf["logerrorfile"] != "" {
		errfile, err := os.OpenFile(conf["logerrorfile"], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Panic(err)
		}
		se = errfile
	}

	logLevel, err := FromString(conf["loglevel"])
	if err != nil {
		log.Fatal(err)
	}
	flags := log.Ldate | log.Ltime | log.Lshortfile

	traceLog := log.New(so, fmt.Sprintf(logFormat, TRACE), flags)
	if logLevel.Index() < TRACE.Index() {
		traceLog = log.New(devNull, "", 0)
	}

	debugLog := log.New(so, fmt.Sprintf(logFormat, DEBUG), flags)
	if logLevel.Index() < DEBUG.Index() {
		debugLog = log.New(devNull, "", 0)
	}

	infoLog := log.New(so, fmt.Sprintf(logFormat, INFO), flags)
	if logLevel.Index() < INFO.Index() {
		infoLog = log.New(devNull, "", 0)
	}

	errorLog := log.New(se, fmt.Sprintf(logFormat, ERROR), flags)
	if logLevel.Index() < ERROR.Index() {
		errorLog = log.New(devNull, "", 0)
	}
	fatalLog := log.New(se, fmt.Sprintf(logFormat, FATAL), flags)

	return &Logger{
		Fatal: fatalLog,
		Error: errorLog,
		Info:  infoLog,
		Debug: debugLog,
		Trace: traceLog,
		Level: logLevel,
	}
}

func init() {
	conf := map[string]string{
		"loglevel":  "info",
		"logformat": "%s\t\t\t\t",
	}
	DefaultLog = NewLogger(conf)
}
