package try_test

import (
	"series"
	"testing"

	cm "github.com/easierway/concurrent_map"
)

func TestPackage(t *testing.T) {
	t.Log(series.GetFibonacciSerie(5))
	t.Log(series.Square(5))
}

// 使用github上第三方包
func TestConcurrentMap(t *testing.T) {
	m := cm.CreateConcurrentMap(99)
	m.Set(cm.StrKey("key"), 20)
	t.Log(m.Get(cm.StrKey("key")))
}

// for的语法
func TestT2(t *testing.T) {
	n := 0
	for n < 5 {
		n++
		t.Log(n)
	}
}

// if的语法
func TestT3(t *testing.T) {
	if a := 2; a == 2 {
		t.Log(a)
	}
}

// switch 语法
func TestT4(t *testing.T) {
	switch a := "d"; a {
	case "a", "b":
		t.Log(a)
	case "c":
		t.Log(a)
	default:
		t.Log("xxxx")
	}
}

func TestT5(t *testing.T) {
	//var arr [3]int
	arr1 := [...]int{1, 2, 3, 4}
	arr2 := [2][2]int{{1, 2}, {3, 4}}
	//idx,e
	for _, e := range arr1 {
		t.Log(e)
	}
	t.Log(arr2)
	t.Log(arr1[2:])
	t.Log(arr1[1:3])
}

func TestSlice(t *testing.T) {
	var s0 []int
	t.Log(len(s0), cap(s0))
	s0 = append(s0, 1)
	s0 = append(s0, 5)
	t.Log(s0)
	t.Log(len(s0), cap(s0))

	s1 := []int{1, 2, 3, 4}
	t.Log(len(s1), cap(s1))

	s2 := make([]int, 3, 5)
	t.Log(len(s2), cap(s2))
}

func TestSliceG(t *testing.T) {
	s := []int{}
	i := 0
	for i < 10 {
		i++
		s = append(s, i)
		t.Log(len(s), cap(s))
	}
}

func TestMapForSet(t *testing.T) {
	mySet := map[int]int{}
	mySet[1] = 1
	n := 3
	if v, ok := mySet[n]; ok {
		t.Logf("%d %d is existing", n, v)
	} else {
		t.Logf("%d %d is not existing", n, v)
	}

	mySet1 := map[int]bool{}
	mySet1[1] = true //如果是false 当成false处理
	n1 := 1
	//if v, ok := mySet1[n1]; ok {
	if mySet1[n1] {
		t.Log(n1, mySet1[n1])
	} else {
		t.Log(false, n1, mySet1[n1])
	}
	delete(mySet1, n1)
	t.Log(len(mySet1), mySet1[n1], mySet1[5])

	mySet3 := map[int]string{}
	mySet3[1] = "aaa"
	mySet3[2] = "bbb"
	n3 := 1
	if v, ok := mySet3[n3]; ok {
		t.Logf("%d %s is existing", n3, v)
	} else {
		t.Logf("%d %s is not existing", n3, v)
	}
	t.Log(mySet3)
	delete(mySet3, 2)
	t.Log(mySet3)
	t.Log("*", mySet3[2], "*")
}
