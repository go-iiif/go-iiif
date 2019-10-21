package tools

import (
	"context"
)

/*

the idea for each and any of the non-native drivers to be able to provide their
own copy of the standard toolset by doing something like this (error handling
omitted for brevity) :

import (
	_ "github.com/go-iiif/go-iiif/native"
	"github.com/go-iiif/go-iiif/tools"
)

func main() {
	tool, _ := tools.NewTransformTool()
	tool.Run()
}

*/

type Tool interface {
	Run(context.Context) error
}
