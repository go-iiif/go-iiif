package log

import (
	"bytes"
	e "errors"
	"fmt"
	"strings"
	"sync"

	"text/template"

	"github.com/go-errors/errors"
	"golang.org/x/net/context"
)

// LogLevel describes a log-level.
type LogLevel int

// Config severity integers.
const (
	// LevelDebug exposes debug logging an above. This exposes all logging.
	LevelDebug LogLevel = iota

	// LevelInfo exposes info logging and above.
	LevelInfo LogLevel = iota

	// LevelWarning exposes warning logging and above.
	LevelWarning LogLevel = iota

	// LevelError exposes error logging and above. This is the most restrictive.
	LevelError LogLevel = iota
)

// Config severity names.

// LogLevelName describes the name of a log-level.
type LogLevelName string

const (
	levelNameDebug   LogLevelName = "debug"
	levelNameInfo    LogLevelName = "info"
	levelNameWarning LogLevelName = "warning"
	levelNameError   LogLevelName = "error"
)

// Seveirty name->integer map.
var (
	levelNameMap = map[LogLevelName]LogLevel{
		levelNameDebug:   LevelDebug,
		levelNameInfo:    LevelInfo,
		levelNameWarning: LevelWarning,
		levelNameError:   LevelError,
	}

	levelNameMapR = map[LogLevel]LogLevelName{
		LevelDebug:   levelNameDebug,
		LevelInfo:    levelNameInfo,
		LevelWarning: levelNameWarning,
		LevelError:   levelNameError,
	}
)

// Other
var (
	includeFilters    = make(map[string]bool)
	useIncludeFilters = false
	excludeFilters    = make(map[string]bool)
	useExcludeFilters = false

	adapters = make(map[string]LogAdapter)

	// TODO(dustin): !! Finish implementing this.
	excludeBypassLevel LogLevel = -1
)

// AddIncludeFilter adds global include filter.
func AddIncludeFilter(noun string) {
	includeFilters[noun] = true
	useIncludeFilters = true
}

// RemoveIncludeFilter removes global include filter.
func RemoveIncludeFilter(noun string) {
	delete(includeFilters, noun)
	if len(includeFilters) == 0 {
		useIncludeFilters = false
	}
}

// AddExcludeFilter adds global exclude filter.
func AddExcludeFilter(noun string) {
	excludeFilters[noun] = true
	useExcludeFilters = true
}

// RemoveExcludeFilter removes global exclude filter.
func RemoveExcludeFilter(noun string) {
	delete(excludeFilters, noun)
	if len(excludeFilters) == 0 {
		useExcludeFilters = false
	}
}

// AddAdapter registers a new adapter.
func AddAdapter(name string, la LogAdapter) {
	if _, found := adapters[name]; found == true {
		Panic(e.New("adapter already registered"))
	}

	if la == nil {
		Panic(e.New("adapter is nil"))
	}

	adapters[name] = la

	if GetDefaultAdapterName() == "" {
		SetDefaultAdapterName(name)
	}
}

// ClearAdapters deregisters all adapters.
func ClearAdapters() {
	adapters = make(map[string]LogAdapter)
	SetDefaultAdapterName("")
}

// LogAdapter describes minimal log-adapter functionality.
type LogAdapter interface {
	// Debugf logs a debug message.
	Debugf(lc *LogContext, message *string) error

	// Infof logs an info message.
	Infof(lc *LogContext, message *string) error

	// Warningf logs a warning message.
	Warningf(lc *LogContext, message *string) error

	// Errorf logs an error message.
	Errorf(lc *LogContext, message *string) error
}

// TODO(dustin): !! Also populate whether we've bypassed an exception so that
//                  we can add a template macro to prefix an exclamation of
//                  some sort.

// MessageContext describes the current logging context and can be used for
// substitution in a template string.
type MessageContext struct {
	Level         *LogLevelName
	Noun          *string
	Message       *string
	ExcludeBypass bool
}

// LogContext encapsulates the current context for passing to the adapter.
type LogContext struct {
	logger *Logger
	ctx    context.Context
}

// Logger is the main logger type.
type Logger struct {
	isConfigured bool
	an           string
	la           LogAdapter
	t            *template.Template
	systemLevel  LogLevel
	noun         string
}

// NewLoggerWithAdapterName initializes a logger struct to log to a specific
// adapter.
func NewLoggerWithAdapterName(noun string, adapterName string) (l *Logger) {
	l = &Logger{
		noun: noun,
		an:   adapterName,
	}

	return l
}

// NewLogger returns a new logger struct.
func NewLogger(noun string) (l *Logger) {
	l = NewLoggerWithAdapterName(noun, "")

	return l
}

// Noun returns the noun that this logger represents.
func (l *Logger) Noun() string {
	return l.noun
}

// Adapter returns the adapter used by this logger struct.
func (l *Logger) Adapter() LogAdapter {
	return l.la
}

var (
	configureMutex sync.Mutex
)

func (l *Logger) doConfigure(force bool) {
	configureMutex.Lock()
	defer configureMutex.Unlock()

	if l.isConfigured == true && force == false {
		return
	}

	if IsConfigurationLoaded() == false {
		Panic(e.New("can not configure because configuration is not loaded"))
	}

	if l.an == "" {
		l.an = GetDefaultAdapterName()
	}

	// If this is empty, then no specific adapter was given or no system
	// default was configured (which implies that no adapters were registered).
	// All of our logging will be skipped.
	if l.an != "" {
		la, found := adapters[l.an]
		if found == false {
			Panic(fmt.Errorf("adapter is not valid: %s", l.an))
		}

		l.la = la
	}

	// Set the level.

	systemLevel, found := levelNameMap[levelName]
	if found == false {
		Panic(fmt.Errorf("log-level not valid: [%s]", levelName))
	}

	l.systemLevel = systemLevel

	// Set the form.

	if format == "" {
		Panic(e.New("format is empty"))
	}

	if t, err := template.New("logItem").Parse(format); err != nil {
		Panic(err)
	} else {
		l.t = t
	}

	l.isConfigured = true
}

func (l *Logger) flattenMessage(lc *MessageContext, format *string, args []interface{}) (string, error) {
	m := fmt.Sprintf(*format, args...)

	lc.Message = &m

	var b bytes.Buffer
	if err := l.t.Execute(&b, *lc); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (l *Logger) allowMessage(noun string, level LogLevel) bool {
	if _, found := includeFilters[noun]; found == true {
		return true
	}

	// If we didn't hit an include filter and we *had* include filters, filter
	// it out.
	if useIncludeFilters == true {
		return false
	}

	if _, found := excludeFilters[noun]; found == true {
		return false
	}

	return true
}

func (l *Logger) makeLogContext(ctx context.Context) *LogContext {
	return &LogContext{
		ctx:    ctx,
		logger: l,
	}
}

type logMethod func(lc *LogContext, message *string) error

func (l *Logger) log(ctx context.Context, level LogLevel, lm logMethod, format string, args []interface{}) error {
	if l.systemLevel > level {
		return nil
	}

	// Preempt the normal filter checks if we can unconditionally allow at a
	// certain level and we've hit that level.
	//
	// Notice that this is only relevant if the system-log level is letting
	// *anything* show logs at the level we came in with.
	canExcludeBypass := level >= excludeBypassLevel && excludeBypassLevel != -1
	didExcludeBypass := false

	n := l.Noun()

	if l.allowMessage(n, level) == false {
		if canExcludeBypass == false {
			return nil
		}

		didExcludeBypass = true
	}

	levelName, found := levelNameMapR[level]
	if found == false {
		Panicf("level not valid: (%d)", level)
	}

	levelName = LogLevelName(strings.ToUpper(string(levelName)))

	mc := &MessageContext{
		Level:         &levelName,
		Noun:          &n,
		ExcludeBypass: didExcludeBypass,
	}

	s, err := l.flattenMessage(mc, &format, args)
	PanicIf(err)

	lc := l.makeLogContext(ctx)

	err = lm(lc, &s)
	PanicIf(err)

	if level == LevelError {
		return e.New(s)
	}

	return nil
}

func (l *Logger) mergeStack(err interface{}, format string, args []interface{}) (string, []interface{}) {
	if format != "" {
		format += "\n%s"
	} else {
		format = "%s"
	}

	var stackified *errors.Error
	stackified, ok := err.(*errors.Error)
	if ok == false {
		stackified = errors.Wrap(err, 2)
	}

	args = append(args, stackified.ErrorStack())

	return format, args
}

// Debugf forwards debug-logging to the underlying adapter.
func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.doConfigure(false)

	if l.la != nil {
		l.log(ctx, LevelDebug, l.la.Debugf, format, args)
	}
}

// Infof forwards debug-logging to the underlying adapter.
func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.doConfigure(false)

	if l.la != nil {
		l.log(ctx, LevelInfo, l.la.Infof, format, args)
	}
}

// Warningf forwards debug-logging to the underlying adapter.
func (l *Logger) Warningf(ctx context.Context, format string, args ...interface{}) {
	l.doConfigure(false)

	if l.la != nil {
		l.log(ctx, LevelWarning, l.la.Warningf, format, args)
	}
}

// Errorf forwards debug-logging to the underlying adapter.
func (l *Logger) Errorf(ctx context.Context, errRaw interface{}, format string, args ...interface{}) {
	l.doConfigure(false)

	var err interface{}

	if errRaw != nil {
		_, ok := errRaw.(*errors.Error)
		if ok == true {
			err = errRaw
		} else {
			err = errors.Wrap(errRaw, 1)
		}
	}

	if l.la != nil {
		if errRaw != nil {
			format, args = l.mergeStack(err, format, args)
		}

		l.log(ctx, LevelError, l.la.Errorf, format, args)
	}
}

// ErrorIff logs a string-substituted message if errRaw is non-nil.
func (l *Logger) ErrorIff(ctx context.Context, errRaw interface{}, format string, args ...interface{}) {
	if errRaw == nil {
		return
	}

	var err interface{}

	_, ok := errRaw.(*errors.Error)
	if ok == true {
		err = errRaw
	} else {
		err = errors.Wrap(errRaw, 1)
	}

	l.Errorf(ctx, err, format, args...)
}

// Panicf logs a string-substituted message.
func (l *Logger) Panicf(ctx context.Context, errRaw interface{}, format string, args ...interface{}) {
	l.doConfigure(false)

	var wrapped interface{}

	_, ok := errRaw.(*errors.Error)
	if ok == true {
		wrapped = errRaw
	} else {
		wrapped = errors.Wrap(errRaw, 1)
	}

	if l.la != nil {
		format, args = l.mergeStack(wrapped, format, args)
		wrapped = l.log(ctx, LevelError, l.la.Errorf, format, args)
	}

	Panic(wrapped)
}

// PanicIff panics with a string-substituted message if errRaw is non-nil.
func (l *Logger) PanicIff(ctx context.Context, errRaw interface{}, format string, args ...interface{}) {
	if errRaw == nil {
		return
	}

	// We wrap the error here rather than rely on on Panicf because there will
	// be one more stack-frame than expected and there'd be no way for that
	// method to know whether it should drop one frame or two.

	var err interface{}

	_, ok := errRaw.(*errors.Error)
	if ok == true {
		err = errRaw
	} else {
		err = errors.Wrap(errRaw, 1)
	}

	l.Panicf(ctx, err.(error), format, args...)
}

// Wrap returns a stack-wrapped error. If already stack-wrapped this is a no-op.
func Wrap(err interface{}) *errors.Error {
	es, ok := err.(*errors.Error)
	if ok == true {
		return es
	}

	return errors.Wrap(err, 1)
}

// Errorf returns a stack-wrapped error with a string-substituted message.
func Errorf(message string, args ...interface{}) *errors.Error {
	err := fmt.Errorf(message, args...)
	return errors.Wrap(err, 1)
}

// Panic panics with the error. Wrap if not already stack-wrapped.
func Panic(err interface{}) {
	_, ok := err.(*errors.Error)
	if ok == true {
		panic(err)
	} else {
		panic(errors.Wrap(err, 1))
	}
}

// Panicf panics a stack-wrapped error with a string-substituted message.
func Panicf(message string, args ...interface{}) {
	err := Errorf(message, args...)
	Panic(err)
}

// PanicIf panics if err is non-nil.
func PanicIf(err interface{}) {
	if err == nil {
		return
	}

	_, ok := err.(*errors.Error)
	if ok == true {
		panic(err)
	} else {
		panic(errors.Wrap(err, 1))
	}
}

// Is checks if the left ("actual") error equals the right ("against") error.
// The right must be an unwrapped error (the kind that you'd initialize as a
// global variable). The left can be a wrapped or unwrapped error.
func Is(actual, against error) bool {
	// If it's an unwrapped error.
	if _, ok := actual.(*errors.Error); ok == false {
		return actual == against
	}

	return errors.Is(actual, against)
}

// PrintError is a utility function to prevent the caller from having to import
// the third-party library.
func PrintError(err error) {
	wrapped := Wrap(err)
	fmt.Printf("Stack:\n\n%s\n", wrapped.ErrorStack())
}

// PrintErrorf is a utility function to prevent the caller from having to
// import the third-party library.
func PrintErrorf(err error, format string, args ...interface{}) {
	wrapped := Wrap(err)

	fmt.Printf(format, args...)
	fmt.Printf("\n")
	fmt.Printf("Stack:\n\n%s\n", wrapped.ErrorStack())
}

func init() {
	if format == "" {
		format = defaultFormat
	}

	if levelName == "" {
		levelName = defaultLevelName
	}

	if includeNouns != "" {
		for _, noun := range strings.Split(includeNouns, ",") {
			AddIncludeFilter(noun)
		}
	}

	if excludeNouns != "" {
		for _, noun := range strings.Split(excludeNouns, ",") {
			AddExcludeFilter(noun)
		}
	}

	if excludeBypassLevelName != "" {
		excludeBypassLevelName = LogLevelName(strings.ToLower(string(excludeBypassLevelName)))

		var found bool
		if excludeBypassLevel, found = levelNameMap[excludeBypassLevelName]; found == false {
			panic(e.New("exclude bypass-level is invalid"))
		}
	}
}
