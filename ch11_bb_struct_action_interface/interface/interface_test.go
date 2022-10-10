package interface_test

import "testing"

type Programmer interface {
	//自定义interface
	WriteHelloWorld() string
}

type GoProgrammer struct {
}

func (g *GoProgrammer) WriteHelloWorld() string {
	return "fmt.Println('Hello World')"
}

func TestClient(t *testing.T) {
	//var p Programmer  p = new(GoProgrammer)
	p := new(GoProgrammer)
	t.Log(p.WriteHelloWorld())
}

//接口变量 Programmer是接口类型
func TestVar(t *testing.T) {
	var prog Programmer = &GoProgrammer{} //类型 数据（实例） 把实现接口的结构体的指针类型 给接口
	t.Log(prog.WriteHelloWorld())
}
