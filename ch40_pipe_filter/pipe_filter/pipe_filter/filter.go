// Package pipefilter is to define the interfaces and the structures for pipe-filter style implementation
package pipefilter

// Request is the input of the filter
type Request interface{}

// Response is the output of the filter
type Response interface{}

// Filter interface is the definition of the data processing components
// Pipe-Filter structure
type Filter interface {
	//接收各种参数 类型不同 所以用空接口
	Process(data Request) (Response, error)
}
