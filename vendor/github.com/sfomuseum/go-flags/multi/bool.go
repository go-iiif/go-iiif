package multi

type MultiBool []bool

func (m *MultiBool) Set(value bool) error {
	*m = append(*m, value)
	return nil
}

func (m *MultiBool) Get() interface{} {
	return *m
}
