package emptystruct

import (
	"fmt"
	"reflect"
	"testing"
)

type EmptyStruct struct{}

type NonEmptyStruct struct {
	Field int
}

func TestFn(t *testing.T) {
	empty := EmptyStruct{}
	nonEmpty := NonEmptyStruct{Field: 10}

	nonEmpty2 := NonEmptyStruct{}

	isEmpty := reflect.DeepEqual(empty, reflect.Zero(reflect.TypeOf(empty)).Interface())
	fmt.Println("empty is empty:", isEmpty) // 输出： empty is empty: false

	isEmpty = reflect.DeepEqual(nonEmpty, reflect.Zero(reflect.TypeOf(nonEmpty)).Interface())
	fmt.Println("nonEmpty is empty:", isEmpty) // 输出： nonEmpty is empty: false

	isEmpty = reflect.DeepEqual(nonEmpty2, reflect.Zero(reflect.TypeOf(nonEmpty2)).Interface())
	fmt.Println("nonEmpty2 is empty:", isEmpty) // 输出： nonEmpty2 is empty: true
}

// 当一个结构体指针为<nil>时，表示该指针并未指向任何有效的结构体实例，即没有分配内存空间给该结构体的实例。这通常发生在以下几种情况：

// 当使用var关键字声明一个结构体指针但没有显式地对其进行初始化时，它会被初始化为<nil>。
// 当使用new函数分配一个结构体指针时，如果没有为其分配实际的结构体实例，那么它也会被初始化为<nil>。

// 在使用一个指向结构体的指针时，需要先确保该指针不为<nil>，否则尝试访问或修改其指向的结构体会导致错误。可以通过检查指针是否为<nil>来避免这类错误。
