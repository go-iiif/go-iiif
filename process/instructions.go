package process

import (
	"context"
	"encoding/json"
	"gocloud.dev/blob"
	"io"
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

func ReadInstructionsFromBucket(ctx context.Context, bucket *blob.Bucket, fname string) (IIIFInstructionSet, error) {

	instructions_fh, err := bucket.NewReader(ctx, fname, nil)

	if err != nil {
		return nil, err
	}

	defer instructions_fh.Close()

	return ReadInstructionsReader(instructions_fh)
}

func ReadInstructions(str_instructions string) (IIIFInstructionSet, error) {

	var raw_instructions []byte

	if strings.HasPrefix(str_instructions, "{") {

		raw_instructions = []byte(str_instructions)
		return ReadInstructionsBytes(raw_instructions)
	}

	path := str_instructions
	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return ReadInstructionsReader(fh)

}

func ReadInstructionsReader(fh io.Reader) (IIIFInstructionSet, error) {

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	return ReadInstructionsBytes(body)
}

func ReadInstructionsBytes(body []byte) (IIIFInstructionSet, error) {

	var instruction_set IIIFInstructionSet

	err := json.Unmarshal(body, &instruction_set)

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
