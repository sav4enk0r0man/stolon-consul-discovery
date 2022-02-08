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
	se := os.Stderr
	so := os.Stdout
	devNull := ioutil.Discard

	format := conf["logformat"]
	logLevel, err := FromString(conf["loglevel"])
	if err != nil {
		log.Fatal(err)
	}
	flags := log.Ldate | log.Ltime | log.Lshortfile

	traceLog := log.New(so, fmt.Sprintf(format, TRACE), flags)
	if logLevel.Index() < TRACE.Index() {
		traceLog = log.New(devNull, "", 0)
	}

	debugLog := log.New(so, fmt.Sprintf(format, DEBUG), flags)
	if logLevel.Index() < DEBUG.Index() {
		debugLog = log.New(devNull, "", 0)
	}

	infoLog := log.New(so, fmt.Sprintf(format, INFO), flags)
	if logLevel.Index() < INFO.Index() {
		infoLog = log.New(devNull, "", 0)
	}

	errorLog := log.New(so, fmt.Sprintf(format, ERROR), flags)
	if logLevel.Index() < ERROR.Index() {
		errorLog = log.New(devNull, "", 0)
	}
	fatalLog := log.New(se, fmt.Sprintf(format, FATAL), flags)

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
