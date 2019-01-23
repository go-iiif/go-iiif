package process

type Processor interface {
	ProcessURIWithInstructionSet(string, IIIFInstructionSet) (map[string]string, error)
	ProcessURIWithInstructions(string, IIIFInstructions) (string, error)
}
