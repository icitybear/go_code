package string_test

import (
	"testing"
)

func TestString(t *testing.T) {
	var s string
	t.Log(s) //初始化为默认零值“”
	s = "hello"
	t.Log(len(s))
	//s[1] = '3' //string是不可变的byte slice
	//s = "\xE4\xB8\xA5" //可以存储任何二进制数据
	s = "\xE4\xBA\xBB\xFF"
	t.Log(s)
	t.Log(len(s))

	s = "中国s"
	t.Log(s)
	t.Log(len(s)) //是byte数
	//访问字符串字符，字符串的内容（纯字节）可以通过标准索引法来获取，在方括号[ ]内写入索引，索引从 0 开始计数（只对纯 ASCII 码的字符串有效）
	//注意：获取字符串中某个字节的地址属于非法行为，例如 &str[i]
	c := []rune(s) //字符串 中 转成 rune unicode码点
	t.Log(c)       //[20013 22269 115]
	t.Log(len(c))
	// unsafe.Sizeof返回变量在内存中占用的字节数(切记，如果是slice，则不会返回这个slice在内存中的实际占用长度)
	//t.Log("rune size:", unsafe.Sizeof(c[0]))
	// 不同编码下
	t.Logf("中 unicode %x", c[0])
	t.Logf("中 UTF8 %x", s)
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
