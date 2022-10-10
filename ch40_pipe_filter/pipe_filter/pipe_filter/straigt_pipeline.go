package pipefilter

// StraightPipeline is composed of the filters, and the filters are piled as a straigt line.
type StraightPipeline struct {
	Name    string
	Filters *[]Filter //多个filter
}

// NewStraightPipeline create a new StraightPipelineWithWallTime
//参数 filters ...Filter 接收多个Filter
func NewStraightPipeline(name string, filters ...Filter) *StraightPipeline {
	return &StraightPipeline{
		Name:    name,
		Filters: &filters,
	}
}

// Process is to process the coming data by the pipeline
func (f *StraightPipeline) Process(data Request) (Response, error) {
	var ret interface{}
	var err error

	for _, filter := range *f.Filters {
		//每个filter进行处理
		ret, err = filter.Process(data)
		if err != nil {
			return ret, err
		}
		//上一个filter的结果 作为下一个filter的参数
		data = ret
	}
	return ret, err
}
