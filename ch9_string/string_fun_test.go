package string_test

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
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
	fmt.Printf("s1 + s2 = %s\n", s3)
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
}

// 字符串拼接
func TestXxx(t *testing.T) {
	// 一次性执行完毕的流程 用+ fmt都无所谓，循环脚本，后台挂起的还是用效率高的
	j1()
	j2()
	j3()
	// 5 > 4 >321
	j5()
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
