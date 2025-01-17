package fn_test

import (
	"fmt"
	"testing"
)

// 延迟执行函数defer 通常用于清理某些资源和处理某些异常
// return之后的语句先执⾏，defer后的语句后执⾏
// 在一个函数中,如果有多个 defer 那它的执行顺序是从后往前执行 入栈出栈，后期对传入参数修改，并不会影响栈内函数的值
var tmpStr string = "hello world"

// panic 程序异常中断,在go里面代表了不可修复的错误,在go中defer在panic之后也是会执行的, 但是注册在panic代码前面的财执行，且值为当时的传入值
func TestDefer(t *testing.T) {
	defer func() {
		fmt.Println("clear resources.", tmpStr)
	}()

	x := 10
	defer func(x int) {
		x++
		jbStr := "city"
		tmpStr = tmpStr + " 修改全局变量tmpStr" //
		fmt.Println(tmpStr, jbStr)

		fmt.Println("defer ", x) // 11
	}(x) //这里调用了x变量 defer后面的函数在入栈的时候保存的是入栈那一刻的值，而当时x的值是10，所以后期对x修改，并不会影响栈内函数的值

	x += 5
	fmt.Println("cur", x) // cur 15
	// x = x/0
	//panic("err") // panic后面的代码不会跑，包括后面才注册的defer
	returnAndDeferFunc()

	fmt.Println(tmpStr) // hello world 正常执行下来的

}

var testInt int

func TestOrder(t *testing.T) {
	testInt := returnAndDeferFunc()
	fmt.Println(testInt) // 输出2
}

func deferFunc() {
	fmt.Println("defer func called...")
	testInt = 1
	fmt.Printf("defer func called...testInt:%d\n", testInt)
}

func returnFunc() int {
	fmt.Println("return func called...")
	testInt = 2
	fmt.Printf("return func called...testInt:%d\n", testInt)
	return testInt
}

// return之后的语句先执⾏（返回值已经记录了，后续defer里再更改不会影响），defer后的语句后执⾏
func returnAndDeferFunc() int {
	defer deferFunc()

	return returnFunc()
}

// return func called...
// return func called...testInt:2    return先接收返回值2 保存下来了
// defer func called...
// defer func called...testInt:1
// 2 return后续的defer 影响不到已保存的返回值
