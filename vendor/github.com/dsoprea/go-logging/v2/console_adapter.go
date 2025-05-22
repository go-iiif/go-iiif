package log

import (
	golog "log"
)

// ConsoleLogAdapter prints logging to STDOUT.
type ConsoleLogAdapter struct {
}

// NewConsoleLogAdapter returns a new ConsoleLogAdapter.
func NewConsoleLogAdapter() LogAdapter {
	return new(ConsoleLogAdapter)
}

// Debugf logs a debugging message.
func (cla *ConsoleLogAdapter) Debugf(lc *LogContext, message *string) error {
	golog.Println(*message)

	return nil
}

// Infof logs an info message.
func (cla *ConsoleLogAdapter) Infof(lc *LogContext, message *string) error {
	golog.Println(*message)

	return nil
}

// Warningf logs a warning message.
func (cla *ConsoleLogAdapter) Warningf(lc *LogContext, message *string) error {
	golog.Println(*message)

	return nil
}

// Errorf logs an error message.
func (cla *ConsoleLogAdapter) Errorf(lc *LogContext, message *string) error {
	golog.Println(*message)

	return nil
}
