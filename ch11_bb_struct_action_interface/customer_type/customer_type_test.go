package customer_type

import (
	"fmt"
	"testing"
	"time"
)

type IntConv func(op int) int //自定义类型

//包一层必包 算函数调用时长
func timeSpent(inner IntConv) IntConv {
	return func(n int) int {
		start := time.Now()
		ret := inner(n)
		fmt.Println(ret)
		fmt.Println("time spent:", time.Since(start).Seconds())
		return ret
	}
}

func slowFun(op int) int {
	fmt.Println(op)
	time.Sleep(time.Second * 1)
	return op + 1000
}

func TestFn(t *testing.T) {
	//包一层
	tsSF := timeSpent(slowFun)
	t.Log(tsSF(10))
}
