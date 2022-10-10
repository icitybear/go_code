package array_test

import "testing"

func TestArrayInit(t *testing.T) {
	var arr [3]int
	arr1 := [4]int{1, 2, 3, 4}
	arr3 := [...]int{1, 3, 4, 5}
	arr1[1] = 5
	//arr3 = append(arr3, 6) 必须是切片才能用append
	// 数组 包括长度和值类型  arr和arr1是不同数据类型 因为长度不一样
	t.Log(arr[1], arr[2])
	t.Log(arr1, arr3)
}

func TestArrayTravel(t *testing.T) {
	arr3 := [...]int{1, 3, 4, 5}
	//arr3[5] = 6 数组初始化固定死个数了 只有切片[]才能动态扩容
	//t.Log(arr3[4])
	for i := 0; i < len(arr3); i++ {
		t.Log(arr3[i])
	}
	for _, e := range arr3 {
		t.Log(e)
	}
}

func TestArraySection(t *testing.T) {
	arr3 := [...]int{1, 2, 3, 4, 5}
	arr3_sec := arr3[:]
	t.Log(arr3_sec)
}
