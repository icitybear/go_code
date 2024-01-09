package slice_test

import (
	"fmt"
	"sort"
	"testing"
)

// sort.Sort 自定义排序
// 首先要自定义类型
// 之后要实现自定义排序  需要实现三个方法  Len 返回长度 Less 比较（升序还是降序） Swap 交换
type arrTest []int

func (arr arrTest) Len() int {
	return len(arr)
}

// < 小于号  —— 升序
// > 大于号 —— 降序
func (arr arrTest) Less(i, j int) bool {
	return (arr)[i] > (arr)[j]
}

func (arr arrTest) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
	return
}

func TestSliceSort(t *testing.T) {

	//a := []int{12, 45, 33, 78, 9, 14}
	a := arrTest{12, 45, 33, 78, 9, 14}
	sort.Sort(a)
	fmt.Println(a)
}

type person struct { //定义对象person
	name string
	age  int
}

type personSlice []person //给[]person绑定对象

// 实现sort包定义的Interface接口
func (s personSlice) Len() int {
	return len(s)
}

// asc
func (s personSlice) Less(i, j int) bool {
	return s[i].age < s[j].age // 使用切片元素结构体的字段
}

func (s personSlice) Swap(i, j int) {
	s[i].age, s[j].age = s[j].age, s[i].age
}

func TestSliceSort2(t *testing.T) {
	p := personSlice{
		person{
			name: "mike",
			age:  13,
		}, person{
			name: "jane",
			age:  12,
		}, person{
			name: "peter",
			age:  14,
		}}
	sort.Sort(p)
	fmt.Println(p) // [{mike 12} {jane 13} {peter 14}]
}

// sort.Slice(slice, func(i, j int) bool)
// 任意类型slice, 比较函数如果省略这个函数，则使用内置的比较函数对切片进行排序
// 1.8引入
func TestSliceSort3(t *testing.T) {
	s := []int{5, 2, 6, 3, 1, 4}
	sort.Slice(s, func(i, j int) bool {
		// < 小于号  —— 升序
		// > 大于号 —— 降序
		return s[i] < s[j]
	})
	fmt.Println(s)
}
