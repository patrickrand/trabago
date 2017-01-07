package trabago

type Pipeline struct {
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) AddJob() *Pipeline {
	return p
}
