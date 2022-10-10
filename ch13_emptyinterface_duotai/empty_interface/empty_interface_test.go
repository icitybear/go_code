package empty_interface

import (
	"fmt"
	"testing"
)

// 空接口标识任意类型 i, ok := p.(int) 转换类型 需要传入接口数据就是该类型
func DoSomething(p interface{}) {

	// if i, ok := p.(int); ok {
	// 	fmt.Println("Integer", i)
	// 	return
	// }
	// if s, ok := p.(string); ok {
	// 	fmt.Println("string", s)
	// 	return
	// }

	// 不是改类型的情况下 就无法转换 i为类型空值(如int是0 string是空字符串) ok为false
	i, ok := p.(int)
	fmt.Println(i, ok)
	//过断⾔来将空接⼝转换为制定类型  类型分支type-switch
	switch v := p.(type) {
	case int:
		fmt.Println("Integer", v)
	case float32:
		fmt.Println("float32", v)
	case float64:
		fmt.Println("float64", v)
	case string:
		fmt.Println("String", v)
	default:
		fmt.Println("Unknow Type", v)
	}
}

func TestEmptyInterfaceAssertion(t *testing.T) {
	DoSomething(10)
	DoSomething(1.444)
	DoSomething("20")
}
