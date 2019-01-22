package log

// Please make this implement the rest of log.Logger
// https://golang.org/pkg/log/

import (
	"fmt"
	"io"
	golog "log"
	"os"
	"path/filepath"
	"strings"
)

type WOFLog interface {
	Fatal(...interface{})
	Error(...interface{})
	Warning(...interface{})
	Status(...interface{})
	Info(...interface{})
	Debug(...interface{})
}

type MockLogger struct{}

func (m *MockLogger) Fatal(...interface{})   {}
func (m *MockLogger) Error(...interface{})   {}
func (m *MockLogger) Warning(...interface{}) {}
func (m *MockLogger) Status(...interface{})  {}
func (m *MockLogger) Info(...interface{})    {}
func (m *MockLogger) Debug(...interface{})   {}

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

func (l WOFLogger) Debug(v ...interface{}) {
	l.dispatch("debug", v...)
}

func (l WOFLogger) Info(v ...interface{}) {
	l.dispatch("info", v...)
}

func (l WOFLogger) Status(v ...interface{}) {
	l.dispatch("status", v...)
}

func (l WOFLogger) Warning(v ...interface{}) {
	l.dispatch("warning", v...)
}

func (l WOFLogger) Error(v ...interface{}) {
	l.dispatch("error", v...)
}

func (l WOFLogger) Fatal(v ...interface{}) {
	l.dispatch("fatal", v...)
	os.Exit(1)
}

func (l WOFLogger) dispatch(level string, args ...interface{}) {

	format := "%v"

	if len(args) > 1 {
		format = args[0].(string)
		args = args[1:]
	}

	for minlevel, logger := range l.Loggers {

		if l.emit(level, minlevel) {

			msg := fmt.Sprintf(format, args...)

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
