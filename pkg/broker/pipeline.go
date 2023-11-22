package broker

type PipelineProcessor interface {
	Do()
}

type Pipeline struct {
	Output
}
