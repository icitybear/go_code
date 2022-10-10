package type_test

import (
	"fmt"
	"math"
	"testing"
)

//类型定义与类型别名的区别 =
type MyInt int64

func TestImplicit(t *testing.T) {
	var a int32 = 1
	var b int64
	b = int64(a)
	var c MyInt
	c = MyInt(b)
	t.Log(a, b, c)
}

func TestYDY(t *testing.T) {
	a := math.MaxInt64
	b := math.MaxFloat64
	c := math.MaxUint32
	t.Log(a, b, c)
}

func TestPoint(t *testing.T) {
	// 指针
	a := 1
	aPtr := &a
	//aPtr = aPtr + 1 // 指针不参与运算
	t.Log(a, aPtr)                      // 指针的值 变量的地址
	fmt.Printf("%d %p", aPtr, aPtr)     // 0xc000014210 的十进制 8246338032801
	fmt.Printf("%d %p %p", a, &a, aPtr) //指针(指针地址) 就是变量a的地址
	t.Logf("%T %T", a, aPtr)            //指针类型 *T
}

func TestString(t *testing.T) {
	// 声明了变量默认初始化零值（该类型的零值：int 为 0，float 为 0.0，bool 为 false，string 为空字符串，
	// 切片、函数、指针变量的默认为 nil 等。所有的内存在 Go 中都是经过初始化的。）
	var a [3]int
	t.Log(a, len(a))
	var s string
	t.Log("*" + s + "*") //初始化零值是“”
	t.Log(len(s))

}
