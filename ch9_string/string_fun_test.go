package string_test

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// 拼接
func TestJoinFn(t *testing.T) {
	s := "A,B,C"
	// s := "" // tag: 注意点
	parts := strings.Split(s, ",")
	if parts == nil {
		t.Log("parts nil") // 不为nil
	}
	if len(parts) == 0 {
		t.Log("parts 0") // 不为0
	}
	fmt.Println(parts, len(parts)) // 结果是[] 1  0空字符串
	// 所以会循环一次
	for key, part := range parts {
		fmt.Println(key)
		t.Log(part)
	}
	t.Log(strings.Join(parts, "-"))
}

// 分割
func TestSplitFn(t *testing.T) {
	s1 := "chihuo@golang"
	arr := strings.Split(s1, "@")
	fmt.Printf("arr is %v\n", arr)
}

// 去掉2边空格
func TestTrimFn(t *testing.T) {
	s1 := " chihuo@golang \n"
	s2 := strings.TrimSpace(s1)
	fmt.Printf("trim space '%s'\n", s2)
}

// 是否含有前缀 后缀
func TestHas(t *testing.T) {
	s1 := "chihuo@golang"
	if strings.HasPrefix(s1, "chihuo") {
		fmt.Printf("%s has prefix chihuo\n", s1)
	}
	if strings.HasPrefix(s1, "@") {
		fmt.Printf("%s has prefix chihuo\n", s1)
	}
	if strings.HasSuffix(s1, "golang") {
		fmt.Printf("%s has suffix golang\n", s1)
	}

	// str := "外拍-沈铭炜-剧情-四月底wh2桃子0501-男用户-穿的很特别-翻剪-镜像二创-他趣"
	str := "web.business.image/202410085d0dbb6ef08fd6a64ab986e"
	fmt.Println(len(str))
	fmt.Println((utf8.RuneCountInString(str)))
}

// 子字符串 直接通过切片下表获取 因为字符船 是不可改变的值类型 又是 []types
func TestSubstr(t *testing.T) {
	s1 := "chihuo@golang"
	s2 := s1[6:len(s1)]
	fmt.Printf("sub string is %s\n", s2)
}

// 转化
func TestConv(t *testing.T) {
	s := strconv.Itoa(10)
	t.Log("str" + s)
	if i, err := strconv.Atoi("10"); err == nil {
		t.Log(10 + i)
	} else {
		t.Log(i)
	}
}

// + 号
func j1() {
	s1 := "chihuo"
	s2 := "golang"
	s3 := s1 + "@" + s2
	fmt.Printf("s1 + s2 = %s\n", s3)
}

// fmt格式化Sprintf
func j2() {
	s1 := "chihuo"
	s2 := "golang"
	s3 := fmt.Sprintf("%s@%s", s1, s2)
	fmt.Printf("s1 + s2 = %s\n", s3)
}

// strings.Join
func j3() {
	s1 := "chihuo"
	s2 := "golang"
	s3 := strings.Join([]string{s1, s2}, "@")
	fmt.Printf("s1 + s2 = %s len:%d\n", s3, len(s3))

	var str []string
	s4 := strings.Join(str, "@")
	fmt.Printf("s4 = %s len:%d\n", s4, len(s4))

	str1 := []string{}
	s5 := strings.Join(str1, "@")
	fmt.Printf("s5 = %s len:%d\n", s5, len(s5))
}

// bytes.Buffer
func j4() {
	var bt bytes.Buffer
	s1 := "chihuo"
	s2 := "golang"
	bt.WriteString(s1)
	bt.WriteString("@")
	bt.WriteString(s2)

	s3 := bt.String() //把Buffer缓存里的转成字符串
	fmt.Printf("s1 + s2 = %s\n", s3)
}

// strings.Builder
func j5() {
	var builder strings.Builder
	s1 := "chihuo"
	s2 := "golang"
	builder.WriteString(s1)
	builder.WriteString("@")
	builder.WriteString(s2)
	s3 := builder.String()
	fmt.Printf("s1 + s2 = %s\n", s3)
	builder.WriteString("citybear apend")
	s4 := builder.String()
	fmt.Println(s4)
}

// 字符串拼接
func TestXxx(t *testing.T) {
	// 一次性执行完毕的流程 用+ fmt都无所谓，循环脚本，后台挂起的还是用效率高的
	// j1()
	// j2()
	j3()
	// 5 > 4 >321
	// j5()
}

// 字字符串查询
func Stripos(haystack string, needle string, offset ...int) int {
	off := 0
	if len(offset) > 0 {
		off = offset[0]
	}
	if off > len(haystack) || off < 0 {
		return -1
	}
	// 全转为小写
	haystack = strings.ToLower(haystack[off:])
	needle = strings.ToLower(needle)
	index := strings.Index(haystack, needle) // strings.Contains(info.Content, actionName)
	if index != -1 {
		return off + index
	}
	return index
}

type Value struct {
	Name  string
	Value int32
}

func TestFmt(t *testing.T) {
	v1 := Value{
		Name:  "val",
		Value: 10,
	}
	s1 := fmt.Sprintf("%d %v %+v %#v %T %p %f", v1.Value, v1, v1, v1, v1, &v1, float64(v1.Value))
	fmt.Printf("format is %s\n", s1)
}

func TestSplit(t *testing.T) {
	// str := "1,"
	str := ""
	tmp := strings.Split(str, ",") // tag:就算返回空数组 长度也是1
	fmt.Printf("tmp is:%+v str_len:%d arr_len:%d \n", tmp, len(str), len(tmp))
	arr := []string{}
	// var arr = make([]string, 1)
	fmt.Printf("arr is:%+v  arr_len:%d \n", arr, len(arr))

	arr1 := []string{}
	// arr1 := []string{"1", "v2", "", "4"}
	jnStr := strings.Join(arr1, ",")
	fmt.Printf("jnStr is:%+v  jnStr_len:%d \n", jnStr, len(jnStr))
}

func TestJoin(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	str := ","
	var strArr []string
	for _, v := range arr {
		strArr = append(strArr, fmt.Sprintf("%d", v))
		// strArr = append(strArr, fmt.Sprintf("'%d'", v)) 特殊写法
	}
	res := strings.Join(strArr, str)
	fmt.Println(res)

	strArr1 := []string{"taqu", "miyou", "4"}
	res1 := strings.Join(strArr1, str)
	fmt.Println(res1)
}

// 在Go语言中，字符串本身并不直接支持修改其内部字符（因为字符串在Go中是不可变的），包括将字符串的首字母小写
func TestUp(t *testing.T) {
	s1 := "ChiHuo"

	fmt.Println(strings.ToLower(s1)) // 全小写

	s2 := "GOLANG_HHH"
	r := []rune(s2)
	if unicode.IsUpper(r[0]) {
		r[0] = unicode.ToLower(r[0])
	}

	fmt.Println(string(r)) // 首字母小写

}
