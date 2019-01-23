package flags

import (
	"errors"
	"fmt"
	"strings"
)

type KeyValueArg struct {
	Key   string
	Value string
}

type KeyValueArgs []*KeyValueArg

func (e *KeyValueArgs) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValueArgs) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, "=")

	if len(kv) != 2 {
		return errors.New("Invalid cache argument")
	}

	a := KeyValueArg{
		Key:   kv[0],
		Value: kv[1],
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValueArgs) ToMap() map[string]string {

	m := make(map[string]string)

	for _, arg := range *e {
		m[arg.Key] = arg.Value
	}

	return m
}
