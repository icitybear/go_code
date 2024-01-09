package array_test

import (
	"fmt"
	"testing"
)

// 数组是固定长度的、同一类型的数据集合
// 数组元素通过 {} 包裹，然后通过逗号分隔多个元素
// 数组与切片 slice区别
func TestArr(t *testing.T) {
	// 需要指定长度和元素类型
	a := [2]int{1, 2}
	b := [...]int{1, 2} // 语法糖省略数组长度的声明,会在编译期自动计算出数组长度。

	c := []int{1, 2}
	t.Logf("a %T %+v", a, a)
	t.Logf("b %T %+v", b, b)
	t.Logf("c %T %+v", c, c)
	c = append(c, 4) // a b时数组 c是切片 append切片才能调用
	t.Logf("c %T %+v", c, c)
	d := b[:1]
	t.Logf("d %T %+v", d, d)
}

func TestArrayEqual(t *testing.T) {
	a := [2]int{1, 2}
	b := [2]int{2, 1}
	if a == b {
		fmt.Println("equal")
	} else {
		fmt.Println("not equal")
	}

}

func TestArrayInit(t *testing.T) {
	// 初始化的时候，如果没有填满，则空位会通过对应的元素类型零值填充
	var arr [3]int
	arr1 := [4]int{1, 2, 3, 4} // 与变量一样，函数内通过 := 进行一次性声明和初始化
	arr3 := [...]int{1, 3, 4, 5}
	arr1[1] = 5
	//arr3 = append(arr3, 6) 必须是切片才能用append
	// 比较时 数组 包括长度和值类型  arr和arr1是不同数据类型 因为长度不一样
	t.Log(arr[1], arr[2])
	t.Log(arr1, arr3)
}

func TestArrayTravel(t *testing.T) {
	// 初始化指定下标位置的元素值
	arr2 := [4]int{1: 34, 2: 5}
	fmt.Println(arr2)
	arr3 := [...]int{1, 3, 4, 5}
	//arr3[5] = 6 数组初始化固定死个数了 只有切片[]才能动态扩容
	//t.Log(arr3[4])
	for i := 0; i < len(arr3); i++ {
		t.Log(arr3[i])
	}
	for _, e := range arr3 {
		t.Log(e)
	}
}

func TestArraySection(t *testing.T) {
	arr3 := [...]int{1, 2, 3, 4, 5}
	arr3_sec := arr3[:]

	arr1 := []int{1, 2, 3, 4, 5}
	var arr2 []int
	arr4 := append(arr2, arr1...)
	t.Log(arr4)
	// first argument to append must be a slice; have untyped nil
	// arr5 := append(nil, arr1...) // 但是不能直接使用nil
	// t.Log(arr5)

	t.Log(arr3_sec)
}

// 数组长度在声明后就不可更改, 长度编译时就能获取，数组的长度是该数组类型的一个内置常量，可以用 Go 语言的内置函数 len() 来获取
// 访问
// 使用数组下标访问 超出这个范围编译时会报索引越界异常 invalid array index 5 (out of bounds for 5-element array)
// 遍历 for len
// 遍历 range range 表达式返回两个值，第一个是数组下标索引值，第二个是索引对应数组元素值

// 多维数组 每个元素可能是个数组，在进行循环遍历的时候需要多层嵌套循环

// 缺点
// 不能动态添加元素到数组
// 值类型，作为参数传递到函数时，传递的是数组的值拷贝，也就是说，会先将数组拷贝给形参，然后在函数体中引用的是形参而不是原来的数组，当我们在函数中对数组元素进行修改时，并不会影响原来的数组

// 需要一个引用类型的、支持动态添加元素的新「数组」类型- 切片类型
