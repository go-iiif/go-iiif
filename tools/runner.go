package tools

import (
	"context"
	"flag"
)

type ToolRunner struct {
	Tool
	tools []Tool
}

func NewToolRunner(tools ...Tool) (Tool, error) {

	t := &ToolRunner{
		tools: tools,
	}

	return t, nil
}

func (t *ToolRunner) Run(ctx context.Context) error {

	fs := flag.NewFlagSet("combined", flag.ExitOnError)
	return t.RunWithFlagSet(ctx, fs)
}

func (t *ToolRunner) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done_ch := make(chan bool)
	err_ch := make(chan error)
	
	for _, tl := range t.tools {

		go func(ctx context.Context, tl Tool, fs *flag.FlagSet) {

			defer func(){
				done_ch <- true
			}()
			
			select {
			case <- ctx.Done():
				return
			default:
				// pass
			}
			
			err := t.RunWithFlagSet(ctx, fs)

			if err != nil {
				err_ch <- err
			}
			
		}(ctx, tl, fs)

	}

	remaining := len(t.tools)

	for remaining > 0 {
		select {
		case <- ctx.Done():
			return nil
		case <- done_ch:
			remaining -= 1
		case err := <- err_ch:
			return err
		}
	}

	return nil
}
