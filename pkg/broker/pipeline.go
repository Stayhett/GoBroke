package broker

type PipelineProcessor interface {
	Do()
	GetData() []byte
	GetOutput() Output
}

type Pipeline struct {
	Output
}
