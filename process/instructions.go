package process

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type IIIFInstructionSet map[Label]IIIFInstructions

type IIIFInstructions struct {
	Region   string `json:"region"`
	Size     string `json:"size"`
	Rotation string `json:"rotation"`
	Quality  string `json:"quality"`
	Format   string `json:"format"`
}

func ReadInstructions(str_instructions string) (IIIFInstructionSet, error) {

	var raw_instructions []byte

	if strings.HasPrefix(str_instructions, "{") {

		raw_instructions = []byte(str_instructions)

	} else {

		path := str_instructions
		fh, err := os.Open(path)

		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return nil, err
		}

		raw_instructions = body
	}

	var instruction_set IIIFInstructionSet

	err := json.Unmarshal(raw_instructions, &instruction_set)

	if err != nil {
		return nil, err
	}

	return instruction_set, nil
}

func EnsureInstructions(i IIIFInstructions) IIIFInstructions {

	if i.Region == "" {
		i.Region = "full"
	}

	if i.Size == "" {
		i.Size = "full"
	}

	if i.Rotation == "" {
		i.Rotation = "0"
	}

	if i.Quality == "" {
		i.Quality = "default"
	}

	if i.Format == "" {
		i.Format = "jpg"
	}

	return i
}
