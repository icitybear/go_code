package pinyin_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mozillazg/go-pinyin"
	"github.com/mozillazg/go-slugify"
)

// github.com/mozillazg/go-pinyin是一个用于将中文字符串转换为拼音的库。它提供了将中文字符串转换为拼音的功能，可以指定拼音的格式（如全拼、首字母等），还可以自定义处理非中文字符的逻辑。主要用于将中文字符串转换为拼音，适用于需要处理中文的场景。
// github.com/mozillazg/go-slugify是一个用于生成URL友好的slug（短标识符）的库。它提供了将字符串转换为slug的功能，会去除特殊字符、空格，并将字符串转换为小写字母，适合用于生成URL中的标识符。主要用于处理字符串，生成易读的URL标识符。
// 总结来说，github.com/mozillazg/go-pinyin主要用于处理中文字符串转换为拼音，而github.com/mozillazg/go-slugify主要用于生成URL友好的slug。根据具体的需求选择合适的库来处理字符串转换和处理

func TestParse(t *testing.T) {
	// 调用函数进行转换
	customFallback := func(r rune, a pinyin.Args) []string {
		return []string{string(r)} // 对于非中文字符串的的处理
	}

	str := "abc你好 啊!-世1界(测试)"
	config := pinyin.Args{
		Style:     pinyin.Normal,
		Heteronym: false,
		Separator: "-",
		Fallback:  customFallback, // 如果为nil 非中文字符串就直接过滤
	}
	pinyinSlice := pinyin.LazyPinyin(str, config)
	for _, item := range pinyinSlice {
		fmt.Println(item)
	}

	// tag: 可以直接拼接 将拼音数组合并为一个字符串, 可以根据特殊符号的字节位置，遇到中文字符再拼接对应的链接符号（字节写）
	pinyinStr := strings.Join(pinyinSlice, "")

	// 将拼音字符串转换为小写并
	pinyinStr = strings.ToLower(pinyinStr)
	pinyinStr = strings.ReplaceAll(pinyinStr, " ", "_")
	fmt.Println(pinyinStr)
}

func Test2(t *testing.T) {
	// 这个命名里-符号会被忽略
	s1 := "abc你好!-世1界(测 试)"
	fmt.Println(slugify.Slugify(s1)) // abcni-hao-shi-1jie-ce-shi

	// 创建一个Slugify对象
	s := slugify.Slugify

	s2 := s("Hello 世界!")
	fmt.Println(s2)
}
