package fn_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 申明函数：
// func 函数名(参数或多个参数 参数类型)(多个返回值的类型){}
// func 函数名()单个返回值的类型{}
// func 函数名(参数 ...参数类型)返回值类型{}
// func(){}() 匿名函数 也能传参

func returnMultiValues() (int, int) {
	return rand.Intn(10), rand.Intn(20)
}

// 函数作为参数 以及函数作为返回值
// timeSpent在这个闭包函数中传入的是一个函数 inner func(op int) int 相当于inner 的类型是func(opt int) int
func timeSpent(inner func(op int) int) func(op int) int {
	return func(n int) int {
		start := time.Now()
		ret := inner(n)
		fmt.Println(ret)
		fmt.Println("time spent:", time.Since(start).Seconds())
		return ret + 10000
	}
}

func slowFun(op int) int {
	fmt.Println("op:", op)
	time.Sleep(time.Second * 3)
	return op + 1000
}

func TestFn(t *testing.T) {
	//_ 忽略另外一个返回值
	a, b := returnMultiValues()
	t.Log(a, b)

	//slowFun 函数变量 赋值给新变量 然后再用来传参 等价于直接写函数名
	// var funcVar func(op int) int
	// funcVar = slowFun          //一般是是会用匿名函数
	// tsSF := timeSpent(funcVar) //slowFun对应的参数和返回值 要符合 timeSpent的声明

	// 函数作为参数 以及函数作为返回值 闭包 10也是slowFun参数
	tsSF := timeSpent(slowFun)
	t.Log(tsSF(10))
}

// 函数 可变参数 被转化成一个数组，然后通过数组遍历
// 可变参数类型约束为 int，如果你希望传任意类型，可以指定类型为 interface{}
func Sum(ops ...int) int {
	ret := 0
	//从内部实现机理上来说，类型...type本质上是一个数组切片，也就是[]type，
	//这也是为什么上面的参数 args 可以用 for 循环来获得每个传入的参数。
	for _, op := range ops {
		ret += op
	}
	return ret
}

// 测试可变参数
func TestVarParam(t *testing.T) {
	t.Log(Sum(1, 2, 3, 4))
	t.Log(Sum(1, 2, 3, 4, 5))
}

// adder 函数返回的时候返回的是个闭包
func adder() func(int) int {
	sum := -1 // 自由变量
	fmt.Printf("adder %d", sum)
	return func(v int) int { // v 是局部变量 这个匿名函数就是一个函数体
		fmt.Printf("+ %d ", v)
		sum += v
		fmt.Printf("sum %d ", sum)
		return sum
	}
}

// 闭包 函数返回的是闭包
func TestBb(t *testing.T) {
	// return func 返回的不是代码,返回是函数以及它的对sum的引用,sum会保存下来,保存到这个函数里面,
	// 初始化后的sum  返回的闭包 会引用这个变量， 传的参数是在变的
	a := adder() //赋值了函数体 函数为返回值
	for i := 0; i < 10; i++ {
		fmt.Println(a(i))
	}
}

// 调用匿名函数
func TestCB(t *testing.T) {
	// 准备一个字符串
	var str string = "hello world"
	t.Log(&str)
	// 创建一个匿名函数
	// var foo func()
	foo := func() {
		// 匿名函数中访问str  &str地址是一样的
		str = "hello dude"
		t.Log(&str)
	}
	t.Log(str)
	foo()
	t.Log(str)
}

// 提供一个值, 每次调用函数会指定对值进行累加 闭包实例 返回闭包
// 函数变量 func() int
func Accumulate(value int) func() int {
	// 返回一个闭包
	return func() int {
		// 累加
		value++
		// 返回一个累加值
		return value
	}
}

func TestBbsl(t *testing.T) {
	// 创建一个累加器, 初始值为1
	accumulator := Accumulate(1)

	// 累加1并打印
	fmt.Println(accumulator()) //2
	fmt.Println(accumulator()) //3

	// 打印累加器的函数地址
	fmt.Printf("%p\n", &accumulator) //0xc00000e048

	// 创建一个累加器, 初始值为10
	accumulator2 := Accumulate(10)
	// 累加1并打印
	fmt.Println(accumulator2()) //11

	// 打印累加器的函数地址
	fmt.Printf("%p\n", &accumulator2) //0xc00000e050
}

// 递归 斐波那契数列
func TestBbsl2(t *testing.T) {
	result := 0
	for i := 1; i <= 10; i++ {
		result = fibonacci(i)
		fmt.Printf("fibonacci(%d) is: %d\n", i, result)
	}
}
func fibonacci(n int) (res int) {
	if n <= 2 {
		res = 1
	} else {
		res = fibonacci(n-1) + fibonacci(n-2)
	}
	return
}

func TestBbsl3(t *testing.T) {
	result := Factorial(3)
	fmt.Printf("Factorial is: %d\n", result)
}

func Factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * Factorial(n-1)
		return result
	}
	return 1
}
