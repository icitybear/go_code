package callback

import (
	"fmt"
	"testing"
)

// 调用器接口
type Invoker interface {
	// 需要实现一个Call方法
	Call(interface{})
}

// 结构体类型
type Struct struct {
}

// 实现Invoker的Call 实现该接口Invoker
func (s *Struct) Call(p interface{}) {
	fmt.Println("from struct", p)
}

// 函数定义为类型（方便面向对象）
type FuncCaller func(interface{})

// 实现Invoker的Call 实现该接口Invoker
func (f FuncCaller) Call(p interface{}) {
	// FuncCaller 的 Call() 方法被调用与 func(interface{}) 无关，还需要手动调用函数本体。
	// 所以这里调用f函数本体
	f(p)
}

func TestDT(t *testing.T) {

	// 声明接口变量
	var invoker Invoker

	// 实例化结构体
	s := new(Struct)

	// 将实例化的结构体赋值到接口
	invoker = s

	// 使用接口调用实例化结构体的方法Struct.Call
	invoker.Call("hello world")

	// 将匿名函数转为FuncCaller类型-类型转换 只要底层一样，再赋值给接口
	invoker = FuncCaller(func(v interface{}) {
		fmt.Println("from function", v)
	})

	// 使用接口调用FuncCaller.Call，内部会调用函数本体
	invoker.Call("hello") // 这里才调用闭包函数 传餐执行
}
