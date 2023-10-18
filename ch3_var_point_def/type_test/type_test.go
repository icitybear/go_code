package type_test

import (
	"fmt"
	"math"
	"testing"
)

// 类型定义与类型别名的区别 类型别名多了=
type MyInt int64

func TestImplicit(t *testing.T) {
	var a int32 = 1
	var b int64
	// 如果是大转小 会存在精度丢失（截断）的情况。
	b = int64(a)
	var c MyInt
	c = MyInt(b)
	t.Log(a, b, c)
}

func TestYDY(t *testing.T) {
	// math 包的常量，默认没有类型，会在引用到的地方自动根据实际类型进行推导
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
	t.Log(a, aPtr)                        // 变量的地址 指针的值
	fmt.Printf("%d %p\n", aPtr, aPtr)     // 0xc000014210 的十进制 8246338032801
	fmt.Printf("%d %p %p\n", a, &a, aPtr) //指针(指针地址) 就是变量a的地址
	t.Logf("%T %T\n", a, aPtr)            //指针类型 *T
}

// 参数是指针 引用传递
func add(n *int) {
	*n++ // 相当于 *n = *n + 1;
	// *ptr  对指针n使用 *操作符，也就是指针取值  *ptr=》指针存的地址取值  指向变量的值
	fmt.Println("add函数结束时：", n, *n)
}

func add2(n int) {
	n++
	fmt.Println("add2函数结束时：", n, &n)
}
func TestPoint2(t *testing.T) {
	var y = 2022
	add2(y)
	fmt.Println("调用add2函数之后：", y, &y)
	// 分清引用传递和值传递
	var yy = &y
	//使用&可以获取某个变量的内存地址
	//用*获取到内存地址所对应的值
	add(yy)
	// 地址是同一个
	fmt.Println("调用add函数之后：", &y, y)
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

func TestVar(t *testing.T) {
	var v *int

	// fmt.Println(*v)
	// 会报错panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	// panic: runtime error: invalid memory address or nil pointer dereference

	fmt.Println(v) //<nil>

	v = new(int)

	fmt.Println(*v) // 0

	fmt.Println(v) //0x1400011a190
}
