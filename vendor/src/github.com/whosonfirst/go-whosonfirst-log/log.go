package log

import (
	"errors"
	"fmt"
	"io"
	golog "log"
	"os"
	"strings"
)

type WOFLog interface {
	Fatal(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Status(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type MockLogger struct{}

func (m *MockLogger) Fatal(format string, v ...interface{})   {}
func (m *MockLogger) Error(format string, v ...interface{})   {}
func (m *MockLogger) Warning(format string, v ...interface{}) {}
func (m *MockLogger) Status(format string, v ...interface{})  {}
func (m *MockLogger) Info(format string, v ...interface{})    {}
func (m *MockLogger) Debug(format string, v ...interface{})   {}

type WOFLogger struct {
	Loggers map[string]*golog.Logger
	levels  map[string]int
	Prefix  string
}

func SimpleWOFLogger(prefix string) *WOFLogger {

	logger := NewWOFLogger(prefix)

	stdout := io.Writer(os.Stdout)
	stderr := io.Writer(os.Stderr)

	logger.AddLogger(stdout, "status")
	logger.AddLogger(stderr, "error")

	return logger
}

func NewWOFLogger(prefix string) *WOFLogger {

	loggers := make(map[string]*golog.Logger)
	levels := make(map[string]int)

	levels["fatal"] = 0
	levels["error"] = 10
	levels["warning"] = 20
	levels["status"] = 25
	levels["info"] = 30
	levels["debug"] = 40

	l := WOFLogger{Loggers: loggers, Prefix: prefix, levels: levels}
	return &l
}

func (l WOFLogger) AddLogger(out io.Writer, minlevel string) (bool, error) {

	_, ok := l.Loggers[minlevel]

	if ok {
		return false, errors.New("log level already defined")
	}

	prefix := l.Prefix

	if !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}

	logger := golog.New(out, prefix, golog.Lmicroseconds)
	l.Loggers[minlevel] = logger

	return true, nil
}

func (l WOFLogger) Debug(format string, v ...interface{}) {
	l.dispatch("debug", format, v...)
}

func (l WOFLogger) Info(format string, v ...interface{}) {
	l.dispatch("info", format, v...)
}

func (l WOFLogger) Status(format string, v ...interface{}) {
	l.dispatch("status", format, v...)
}

func (l WOFLogger) Warning(format string, v ...interface{}) {
	l.dispatch("warning", format, v...)
}

func (l WOFLogger) Error(format string, v ...interface{}) {
	l.dispatch("error", format, v...)
}

func (l WOFLogger) Fatal(format string, v ...interface{}) {
	l.dispatch("fatal", format, v...)
	os.Exit(1)
}

func (l WOFLogger) dispatch(level string, format string, v ...interface{}) {

	for minlevel, logger := range l.Loggers {

		if l.emit(level, minlevel) {

			msg := fmt.Sprintf(format, v...)
			out := fmt.Sprintf("[%s] %s", level, msg)
			logger.Println(out)
		}
	}
}

func (l WOFLogger) emit(level string, minlevel string) bool {

	this_level, ok := l.levels[level]

	if !ok {
		return false
	}

	min_level, ok := l.levels[minlevel]

	if !ok {
		return false
	}

	if this_level > min_level {
		return false
	}

	return true
}
