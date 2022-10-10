package string_test

import (
	"strconv"
	"strings"
	"testing"
)

func TestStringFn(t *testing.T) {
	s := "A,B,C"
	parts := strings.Split(s, ",")
	for _, part := range parts {
		t.Log(part)
	}
	t.Log(strings.Join(parts, "-"))
}

func TestConv(t *testing.T) {
	s := strconv.Itoa(10)
	t.Log("str" + s)
	if i, err := strconv.Atoi("10"); err == nil {
		t.Log(10 + i)
	} else {
		t.Log(i)
	}
}

func TestDd(t *testing.T) {
	s := map[string]int{}
	s["a"] = 0
	//ok 是布尔类型
	if v, ok := s["a"]; ok {
		t.Log(v, ok)
	}

	//数组切片类型
	s1 := []int{}
	s1 = append(s1, 1)
	t.Log(s1)

	//v, ok := s1[0]; ok  只有map才是有2个返回值 ok为对应元素是否存在
	// 还不能使用s1[0]值必须使用s1切片 != nil nil是一个预先声明的标识符，只能表示指针、通道、函数、接口、映射或切片
	if s1 != nil {
		t.Log(s1)
	}

	var s2 []int
	//s2 := []int{}

	//直接用现有s1
	s2 = append(s1, 5)
	// [1, 5]
	t.Log(s2)
	for idx, elem := range s2 {
		t.Log(idx, elem)
	}
	s3 := s2[1:]
	t.Log(s3)
}
