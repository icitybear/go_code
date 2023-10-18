package base_var //尽量与文件目录名一致
// 变量
import (
	"fmt"
	"testing"
	"unicode/utf8"
)

// 全是单元测试的方法 并不能go run
func TestVar(t *testing.T) {
	// 变量命名规则遵循驼峰命名法，即首个单词小写，每个新单词的首字母大写，如 userName，但如果你的全局变量希望能够被外部包所使用，则需要将首个单词的首字母也大写。
	// 只声明，不赋值，默认是0 变量在声明之后，系统会自动将变量值初始化为对应类型的零值
	// int 为 0，float 为 0.0，bool 为 false，string 为空字符串，切片、函数、指针变量的默认为 nil 等。所有的内存在 Go 中都是经过初始化的
	var v0 int
	// 先声明 再赋值
	var v1 int
	v1 = 1
	// 声明的同时赋值
	var v2 int = 2
	// 没声明类型 自动推断
	var v3 = 3
	// 短变量声明赋值 只能用在函数内部的写法
	// 多重赋值变量的左值和右值按从左到右的顺序赋值, 配合匿名变量_
	// 推导声明写法的左值变量必须是没有定义过的变量。若定义过，将会发生编译错误。
	v4 := 4

	fmt.Printf("v0=%v, v1=%v, v2=%v, v3=%v, v4=%v\n", v0, v1, v2, v3, v4)

	// 批量声明
	var (
		a1 int
		a2 int = 2
		a3     = 3
	)

	fmt.Printf("a1=%v, a2=%v, a3=%v", a1, a2, a3)
}

func TestInt(t *testing.T) {
	var b uint8 = 0b00000001 // 注意整型范围 超过的话会提示错误
	fmt.Printf("%v, %T\n", b, b)
	// 没声明类型 自动推断 会默认更大范围的类型
	var c = 0b00000001
	fmt.Printf("%v, %T\n", c, c)
	// 变量运算, 类型不一致 会提示int与uint8不匹配 不能运算，所以加强转
	c = 256
	d := uint8(c) / b //c强转时如果范围大于类型值会截取低8bit都是0 256就是0
	// d := c / int(b)
	fmt.Printf("%v, %T\n", d, d)
}

func TestChar(t *testing.T) {

	// byte uint8的别名   1个字节的字符
	// rune int32的别名 4个字节的字符
	var c0 uint8 = 65
	var c1 byte = 65
	if c0 == c1 {
		fmt.Println("byte==uint8")
	}
	//单引号表示单字符 双引号标识字符串 因为中文是多字节的所以必须用rune 4字节
	var c2 = 'a'
	var c3 rune = '中'
	var c4 int32 = 20013 //如果是非int32,比如uint32就不等
	if c3 == c4 {
		fmt.Println("rune==int32")
	}
	fmt.Printf("c0 的码值=%v, 字符为：%c, 类型是%T\n", c0, c0, c0)
	fmt.Printf("c1 的码值=%v, 字符为：%c, 类型是%T\n", c1, c1, c1)
	fmt.Printf("c2 的码值=%v, 字符为：%c, 类型是%T\n", c2, c2, c2)
	fmt.Printf("c3 的码值=%v, 字符为：%c, 类型是%T\n", c3, c3, c3)
}

func TestString(t *testing.T) {
	//双引号表示一个字符串，双引号内字符可以转义
	fmt.Println("\"zifu\tchuan\"")
	// 单引号 单引号只能用来包裹一个字节的ASCII码字符byte 也可以是多字节的字符 rune
	str := 'z'
	fmt.Println(str)             // 输出122
	fmt.Println(`"zifu\tchuan"`) //反引号引起来的字符串就不支持转义
	fmt.Println("'zifuchuan'")

	// 字符串拼接和访问其他练习记录
	// len字符串 字节数 理解字符串Unicode（UTF-8），ASCII字符集
	theme := "中国\ta bc"
	l := len(theme) //\t和空格各算一个字节 中文3个字节 输出11
	fmt.Println(l)
	fmt.Println(theme[7]) //a 97
	for i := 0; i < l; i++ {
		fmt.Printf("ascii: %c  %d\n", theme[i], theme[i])
	}
	l = utf8.RuneCountInString(theme) //7个utf8字符
	fmt.Println(l)
	for _, s := range theme {
		fmt.Printf("Unicode: %c  %d\n", s, s)
	}
}

func TestStringByteRune(t *testing.T) {
	s0 := "中国\ta bc"
	fmt.Printf("值=%v, 类型是%T\n", s0, s0)
	s1 := []rune(s0) //字符串 中 转成 rune unicode码点
	fmt.Printf("值=%v, 类型是%T\n", s1, s1)
	s2 := []byte(s0) //字符串 中 转成 byte字节切片
	fmt.Printf("值=%v, 类型是%T\n", s2, s2)
	// 遍历切片
	for _, s := range s1 {
		fmt.Printf("int32: %c  %d\n", s, s)
	}
	for _, s := range s2 {
		fmt.Printf("uint8: %c  %d\n", s, s)
	}
}

func TestFibList(t *testing.T) {

	a := 1
	b := 1
	t.Log(a)
	for i := 0; i < 5; i++ {
		t.Log(" ", b)
		tmp := a
		a = b
		b = tmp + a
	}

}

func TestExchange(t *testing.T) {

	a := 1
	b := 2
	// tmp := a
	// a = b
	// b = tmp
	// 多重赋值 交换变量
	a, b = b, a
	t.Log(a, b)
}

func TestPtr(t *testing.T) {
	// 指针
	ptr := new(string)
	*ptr = "Go语言教程"
	t.Log(*ptr)
}
