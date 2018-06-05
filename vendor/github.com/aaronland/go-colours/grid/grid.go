package grid

import (
	"errors"
	"github.com/aaronland/go-colours"
	"strings"
)

func NewNamedGrid(name string, args ...interface{}) (colours.Grid, error) {

	var gr colours.Grid
	var err error

	switch strings.ToUpper(name) {

	case "EUCLIDIAN":
		gr, err = NewEuclidianGrid(args)
	default:
		err = errors.New("Unknown or invalid grid")
	}

	return gr, err
}
