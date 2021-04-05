package lookup

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
)

func Lookup(fl *flag.FlagSet, k string) (interface{}, error) {

	v := fl.Lookup(k)

	if v == nil {
		msg := fmt.Sprintf("Unknown flag '%s'", k)
		return nil, errors.New(msg)
	}

	// Go is weird...
	return v.Value.(flag.Getter).Get(), nil
}

func MultiStringVar(fl *flag.FlagSet, k string) ([]string, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return nil, err
	}

	return i.(multi.MultiString), nil
}

func StringVar(fl *flag.FlagSet, k string) (string, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return "", err
	}

	return i.(string), nil
}

func MultiIntVar(fl *flag.FlagSet, k string) ([]int, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return nil, err
	}

	return i.(multi.MultiInt), nil
}

func IntVar(fl *flag.FlagSet, k string) (int, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return 0, err
	}

	return i.(int), nil
}

func MultiInt64Var(fl *flag.FlagSet, k string) ([]int64, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return nil, err
	}

	return i.(multi.MultiInt64), nil
}

func Int64Var(fl *flag.FlagSet, k string) (int64, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return 0, err
	}

	return i.(int64), nil
}

func Float64Var(fl *flag.FlagSet, k string) (float64, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return 0, err
	}

	return i.(float64), nil
}

func BoolVar(fl *flag.FlagSet, k string) (bool, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return false, err
	}

	return i.(bool), nil
}
