package tools

import (
	"context"
	"flag"
)

type ToolRunner struct {
	Tool
	tools    []Tool
	throttle chan bool
}

func NewToolRunner(tools ...Tool) (Tool, error) {

	t := &ToolRunner{
		tools: tools,
	}

	return t, nil
}

func NewSynchronousToolRunner(tools ...Tool) (Tool, error) {

	throttle := make(chan bool)

	t := &ToolRunner{
		tools:    tools,
		throttle: throttle,
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

			err := tl.RunWithFlagSetAndPaths(ctx, fs, paths...)

			if err != nil {
				err_ch <- err
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

	return nil
}
