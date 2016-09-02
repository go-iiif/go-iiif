package source

import ()

type MemorySource struct {
	Source
	body []byte
}

func NewMemorySource(body []byte) (*MemorySource, error) {

	mem := MemorySource{
		body: body,
	}

	return &mem, nil
}

func (mem *MemorySource) Read(uri string) ([]byte, error) {

	return mem.body, nil
}
