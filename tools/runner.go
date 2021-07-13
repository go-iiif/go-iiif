package tools

import (
	"context"
	"errors"
	"flag"
	"fmt"
	_ "log"
)

type ToolRunner struct {
	Tool
	tools       []Tool
	throttle    chan bool
	on_complete ToolRunnerOnCompleteFunc
}

type ToolRunnerOnCompleteFunc func(context.Context, string) error

type ToolRunnerOptions struct {
	Tools          []Tool
	Throttle       chan bool
	OnCompleteFunc ToolRunnerOnCompleteFunc
}

func NewToolRunnerThrottle() (chan bool, error) {

	throttle := make(chan bool)
	return throttle, nil
}

func NewToolRunner(tools ...Tool) (Tool, error) {

	opts := &ToolRunnerOptions{
		Tools: tools,
	}

	return NewToolRunnerWithOptions(opts)
}

func NewSynchronousToolRunner(tools ...Tool) (Tool, error) {

	throttle, err := NewToolRunnerThrottle()

	if err != nil {
		return nil, err
	}

	opts := &ToolRunnerOptions{
		Tools:    tools,
		Throttle: throttle,
	}

	return NewToolRunnerWithOptions(opts)
}

func NewToolRunnerWithOptions(opts *ToolRunnerOptions) (Tool, error) {

	if len(opts.Tools) == 0 {
		return nil, errors.New("No tools to run!")
	}

	t := &ToolRunner{
		tools:       opts.Tools,
		throttle:    opts.Throttle,
		on_complete: opts.OnCompleteFunc,
	}

	return t, nil
}

func (t *ToolRunner) Run(ctx context.Context) error {

	fs := flag.NewFlagSet("combined", flag.ExitOnError)
	return t.RunWithFlagSet(ctx, fs)
}

func (t *ToolRunner) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	paths := fs.Args()
	return t.RunWithFlagSetAndPaths(ctx, fs, paths...)
}

func (t *ToolRunner) RunWithFlagSetAndPaths(ctx context.Context, fs *flag.FlagSet, paths ...string) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, path := range paths {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		done_ch := make(chan bool)
		err_ch := make(chan error)

		if t.throttle != nil {
			go func() {
				t.throttle <- true
			}()
		}

		for _, tl := range t.tools {

			go func(ctx context.Context, tl Tool, fs *flag.FlagSet) {

				if t.throttle != nil {

					<-t.throttle

					defer func() {
						t.throttle <- true
					}()

				}

				defer func() {
					done_ch <- true
				}()

				select {
				case <-ctx.Done():
					return
				default:
					// pass
				}

				err := tl.RunWithFlagSetAndPaths(ctx, fs, path)

				if err != nil {
					err_ch <- fmt.Errorf("Failed to run tool '%T': %w", tl, err)
				}

			}(ctx, tl, fs)

		}

		remaining := len(t.tools)

		for remaining > 0 {
			select {
			case <-ctx.Done():
				return nil
			case <-done_ch:
				remaining -= 1
			case err := <-err_ch:
				return err
			}
		}

		if t.on_complete != nil {

			err := t.on_complete(ctx, path)

			if err != nil {
				return fmt.Errorf("OnComplete function for '%s' failed, %w", path, err)
			}
		}
	}

	return nil
}
