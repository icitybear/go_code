package constant_test

import "testing"

const (
	Monday = 1 + iota //+1 iota一开始默认0
	Tuesday
	Wednesday
)

const (
	Readable = 1 << iota //左移位
	Writable
	Executable
)

func TestConstantTry(t *testing.T) {
	t.Log(Monday, Tuesday)
}

func TestConstantTry1(t *testing.T) {
	t.Log(Readable, Writable, Executable)
	a := 1 //0001
	t.Log(a&Readable == Readable, a&Writable == Writable, a&Executable == Executable)
}
