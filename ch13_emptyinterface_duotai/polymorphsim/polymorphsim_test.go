package polymorphism

import (
	"fmt"
	"testing"
)

type Code string

// 类对接口的实现 多态的使用场景
type Programmer interface {
	//接口要实现的方法 struct与interface
	WriteHelloWorld() Code
}

type GoProgrammer struct {
}

func (p *GoProgrammer) WriteHelloWorld() Code {
	return "fmt.Println(\"Hello World!\")"
}

type JavaProgrammer struct {
}

// 如果 绑定的不是指针是实例 那么传参更方便 但是会有拷贝内存的消耗
func (p *JavaProgrammer) WriteHelloWorld() Code {
	return "System.out.Println(\"Hello World!\")"
}

// 参数要求Programmer接口 以上都实现了接口 指针实例和值实例都实现了
func writeFirstProgram(p Programmer) {
	// WriteHelloWorld要求p是指针 是指针方法
	// 这里的左值可寻址，自动取地址& 调用（形参不算）
	// 多态 根据p的类型  子类 实现了接口
	fmt.Printf("%T %v\n", p, p.WriteHelloWorld())
}

func TestPolymorphism(t *testing.T) {
	goProg := &GoProgrammer{}
	javaProg := new(JavaProgrammer) //可以用&实例来代替 new相当于

	// 如果不用结构体的指针类型会报错 goProg := GoProgrammer{} 因为绑定的方法作用对象是p *GoProgrammer
	// cannot use goProg (variable of type GoProgrammer) as Programmer value in argument to writeFirstProgram:
	// GoProgrammer does not implement Programmer (method WriteHelloWorld has pointer

	writeFirstProgram(goProg)
	writeFirstProgram(javaProg)
}
