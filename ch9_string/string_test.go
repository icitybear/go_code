package string_test

import (
	"fmt"
	"testing"
	"time"
	"unicode/utf8"
	"unsafe"
)

// 字符串

func TestBase(t *testing.T) {
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
	for _, s := range s2 {
		fmt.Printf("uint8: %c  %d\n", s, s)
	}
}

func TestString(t *testing.T) {
	var s string
	t.Log(s) //初始化为默认零值“”
	s = "hello"
	t.Log(len(s))
	//s[1] = '3' //string是不可变的byte slice 只能访问
	//s = "\xE4\xB8\xA5" //可以存储任何二进制数据
	s = "\xE4\xBA\xBB\xFF"
	t.Log(s)
	t.Log(len(s))

	s = "中国s"
	t.Log(s)
	t.Log(len(s)) //是byte数 如果算中文个数用utf8  => 7
	//访问字符串字符，字符串的内容（纯字节）可以通过标准索引法来获取，
	//在方括号[ ]内写入索引，索引从 0 开始计数（只对纯 ASCII 码的字符串有效）
	//注意：获取字符串中某个字节的地址属于非法行为，例如 &str[i]

	c := []rune(s) //字符串 中 转成 rune unicode码点  utf8编码
	t.Log(c)       //[20013 22269 115]
	t.Log(len(c))  //这里也就是c的切片长度了 =>3

	// unsafe.Sizeof返回变量在内存中占用的字节数(切记，如果是slice，则不会返回这个slice在内存中的实际占用长度)
	t.Log("rune size:", unsafe.Sizeof(c[0])) // =>4
	// 不同编码下
	t.Logf("中 unicode %x %v %c", c[0], c[0], c[0]) // unicode 4e2d 20013 中
	t.Logf("中 UTF8 %x %v", s, s)                   //字符串地址 e4b8ade59bbd73
}

func TestStringToRune(t *testing.T) {
	s := "中华人民共和国"
	//byte[] UNICODE
	for _, c := range s {
		//unicode码点  转相应Unicode码点所表示的字符 x 16进制
		// [1] 代表使用参数c
		// t.Logf("%[1]c %[1]x", c)
		t.Logf("%v %c %x", c, c, c)
	}
}

func TestTime(t *testing.T) {

	ti, err := time.Parse("2006-01-02", "2012-03-12")
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return
	}

	formattedString := ti.Format("20060102") // 所需的字符串格式
	fmt.Println(formattedString)
}

func TestTime2(t *testing.T) {
	// tr, err := ConvertStringToTime("2012-03-12 12:20:10")
	tr, err := ConvertStringToTime("2012-03-12")
	if err != nil {
		fmt.Println("日期解析错误:", err)
	}
	fmt.Println(tr)
}

func ConvertStringToTime(str string) (time.Time, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func Test2(t *testing.T) {
	// 默认使用UTC
	endTimer, _ := time.Parse("2006-01-02", "2023-11-23")
	// 2023-11-23 00:00:00 +0000 UTC 但是数据库存放时会加入自己的时区
	// 存入数据库datetime时 存的是本地时间 2023-11-23 08:00:00
	fmt.Print(endTimer)       // 2023-11-23 00:00:00 +0000 UTC
	fmt.Print(endTimer.UTC()) // 2023-11-23 00:00:00 +0000 UTC

	// go指定使用当地时区
	endTimer2, _ := time.ParseInLocation("2006-01-02", "2023-11-23", time.Local)
	fmt.Print(endTimer2) // 2023-11-23 00:00:00 +0800 CST
	// 当地时区时间转成UTC时间
	fmt.Print(endTimer2.UTC()) // 2023-11-22 16:00:00 +0000 UTC
}
