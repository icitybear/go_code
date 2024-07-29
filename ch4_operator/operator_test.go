package operator_test

import "testing"

const (
	Readable = 1 << iota
	Writable
	Executable
)

// 在 Go 语言中，也支持自增/自减运算符，即 ++/--，但是只能作为语句，不能作为表达式，<font color="red">且只能用作后缀，不能放到变量前面</font>, 支持快捷写法
// 位运算 &^ 按位置零 & | ^(异或与取反) << >>
// 逻辑运算 && || ! (只放在bool)
// 比较运算符会考虑变量的类型，<font color="red">各种类型的整型变量都可以直接与字面常量进行比较</font>, >、<、==、>=、<= 和 != 运算结果是布尔值

// 比较
func TestCompareArray(t *testing.T) {
	a := [...]int{1, 2, 3, 4}
	b := [...]int{1, 3, 2, 4}
	//	c := [...]int{1, 2, 3, 4, 5}
	d := [...]int{1, 2, 3, 4}
	t.Log(a == b)
	//t.Log(a == c)
	t.Log(a == d)
}

func TestBitClear(t *testing.T) {
	//&^ 按位置零
	a := 7 //0111
	a = a &^ Readable
	a = a &^ Executable
	t.Log(a)
	// a 为 0010
	t.Log(a&Readable == Readable, a&Writable == Writable, a&Executable == Executable)
}

func TestJisuan(t *testing.T) {

	a := 710 / 100 // 与php不一样会 吞掉小数点 只保留int
	t.Log(a)
	b := float32(710) / float32(100)
	t.Log(b)
}
