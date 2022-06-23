# go-roster

Go package provides interfaces and methods for defining internal lookup tables (or "rosters") for registering and instantiatinge custom interfaces with multiple implementations.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-roster.svg)](https://pkg.go.dev/github.com/aaronland/go-roster)

## Example

The following example is the body of the [roster_test.go](roster_test.go) file:

```
package roster

import (
	"context"
	"fmt"
	"net/url"
	"testing"
)

// Create a toy interface that might have multiple implementations including a common
// method signature for creating instantiations of that interface.

type Example interface {
	String() string
}

type ExampleInitializationFunc func(context.Context, string) (Example, error)

func RegisterExample(ctx context.Context, scheme string, init_func ExampleInitializationFunc) error {
	return example_roster.Register(ctx, scheme, init_func)
}

func NewExample(ctx context.Context, uri string) (Example, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	scheme := u.Scheme

	i, err := example_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to find registeration for %s, %w", scheme, err)
	}

	init_func := i.(ExampleInitializationFunc)
	return init_func(ctx, uri)
}

// Something that implements the Example interface

type StringExample struct {
	Example
	value string
}

func NewStringExample(ctx context.Context, uri string) (Example, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	s := &StringExample{
		value: u.Path,
	}

	return s, nil
}

func (e *StringExample) String() string {
	return e.value
}

// Create a global "roster" of implementations of the Example interface

var example_roster Roster

// Ensure that there is a valid roster (for use by the code handling the Example interface)
// and register the StringExample implementation

func init() {

	ctx := context.Background()

	r, err := NewDefaultRoster()

	if err != nil {
		panic(err)
	}

	example_roster = r

	err = RegisterExample(ctx, "string", NewStringExample)

	if err != nil {
		panic(err)
	}
}

func TestRoster(t *testing.T) {

	ctx := context.Background()

	e, err := NewExample(ctx, "string:///helloworld")

	if err != nil {
		t.Fatalf("Failed to create new example, %v", err)
	}

	v := e.String()

	if v != "/helloworld" {
		t.Fatalf("Unexpected result: '%s'", v)
	}
}
```

## Concrete examples

* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-writer