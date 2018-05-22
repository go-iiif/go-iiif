package log

import (
	"fmt"
	"io"
	golog "log"
	"os"
	"path/filepath"
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
	writers map[string][]io.Writer
	Prefix  string
}

func Prefix(args ...string) string {

	whoami := os.Args[0]
	whoami = filepath.Base(whoami)

	prefix := fmt.Sprintf("[%s]", whoami)

	for _, s := range args {
		prefix = fmt.Sprintf("%s[%s]", prefix, s)
	}

	return prefix
}

func SimpleWOFLogger(args ...string) *WOFLogger {

	logger := NewWOFLogger(args...)

	stderr := io.Writer(os.Stderr)
	logger.AddLogger(stderr, "error")

	// stdout := io.Writer(os.Stdout)
	// logger.AddLogger(stdout, "status")

	return logger
}

func NewWOFLogger(args ...string) *WOFLogger {

	prefix := Prefix(args...)

	writers := make(map[string][]io.Writer)

	loggers := make(map[string]*golog.Logger)
	levels := make(map[string]int)

	levels["fatal"] = 0
	levels["error"] = 10
	levels["warning"] = 20
	levels["status"] = 25
	levels["info"] = 30
	levels["debug"] = 40

	l := WOFLogger{
		Loggers: loggers,
		Prefix:  prefix,
		levels:  levels,
		writers: writers,
	}

	return &l
}

func (l WOFLogger) AddLogger(out io.Writer, minlevel string) (bool, error) {

	_, ok := l.writers[minlevel]

	if !ok {
		wr := make([]io.Writer, 0)
		l.writers[minlevel] = wr
	}

	// check to see whether we already have this logger?

	l.writers[minlevel] = append(l.writers[minlevel], out)

	m := io.MultiWriter(l.writers[minlevel]...)

	logger := golog.New(m, "", golog.Lmicroseconds)
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

			out := fmt.Sprintf("%s %s %s", l.Prefix, strings.ToUpper(level), msg)
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
