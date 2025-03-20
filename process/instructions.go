package process

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aaronland/gocloud-blob/bucket"
	iiifdefaults "github.com/go-iiif/go-iiif/v6/defaults"
	"gocloud.dev/blob"
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

	instructions_r, err := bucket.NewReader(ctx, fname, nil)

	if err != nil {
		return nil, err
	}

	defer instructions_r.Close()

	return ReadInstructionsReader(instructions_r)
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

func ReadInstructionsReader(r io.Reader) (IIIFInstructionSet, error) {

	body, err := io.ReadAll(r)

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

func LoadInstructions(ctx context.Context, bucket_uri string, key string) (IIIFInstructionSet, error) {

	if bucket_uri == iiifdefaults.URI {

		key = "instructions.json"

		r, err := iiifdefaults.FS.Open(key)

		if err != nil {
			return nil, fmt.Errorf("Failed to load config (%s) from defaults, %w", key, err)
		}

		return ReadInstructionsReader(r)
	}

	instructions_bucket, err := bucket.OpenBucket(ctx, bucket_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open instructions bucket, %w", err)
	}

	defer instructions_bucket.Close()

	return ReadInstructionsFromBucket(ctx, instructions_bucket, key)
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
