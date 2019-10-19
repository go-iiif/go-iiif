package session

import (
	"fmt"
)

type AssignEnvError struct {
	err   error
	key   string
	value string
}

func (e *AssignEnvError) String() string {
	return e.Error()
}

func (e *AssignEnvError) Error() string {
	return fmt.Sprintf("Can not assign %s=%s (%s)", e.key, e.value, e.err)
}
