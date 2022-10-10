package slice_test

import "testing"

func TestSliceInit(t *testing.T) {
	var s0 []int
	t.Log(len(s0), cap(s0))
	s0 = append(s0, 1) //开头追加元素
	// 扩容在1024以下 一般是双倍
	t.Log(len(s0), cap(s0))

	s1 := []int{1, 2, 3, 4}
	t.Log(len(s1), cap(s1))

	s2 := make([]int, 3, 5)
	t.Log(len(s2), cap(s2))
	t.Log(s2[0], s2[1], s2[2])
	s2 = append(s2, 1)
	t.Log(s2[0], s2[1], s2[2], s2[3])
	t.Log(len(s2), cap(s2))
}

func TestSliceGrowing(t *testing.T) {
	s := []int{}
	for i := 0; i < 10; i++ {
		s = append(s, i)
		t.Log(len(s), cap(s))
	}
}

func TestSliceShareMemory(t *testing.T) {
	//下标是0开始计数
	year := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep",
		"Oct", "Nov", "Dec"}
	t.Log(year, len(year), cap(year))
	//slice[] （定位到开始位置，然后切对应大小）
	//不改变原来的内存地址 切片包含 开始位置 长度 容量大小
	Q3 := year[8:12] //超出长度会编译报错 不包括结尾位置
	t.Log(Q3, len(Q3), cap(Q3))
	Q4 := year[:0] //缺省写法 [] 0 12
	t.Log(Q4, len(Q4), cap(Q4))
	//year = append(year, "end")
	// 容量问题
	Q2 := year[3:6] //不包括结尾位置
	t.Log(Q2, len(Q2), cap(Q2))
	summer := year[5:8]
	t.Log(summer, len(summer), cap(summer))
	//更改其中一个的值后 year也跟着变了
	summer[0] = "Unknow" // year[5] Q2[2]
	t.Log(Q2)
	t.Log(year)
}

func TestSliceComparing(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{1, 2, 3, 4} //切片
	//c := [...]int{1, 2, 3, 4} 数组
	// if a == b { //切片只能和nil比较
	// 	t.Log("equal")
	// }
	t.Log(a, b)
	var c []int        //只声明 不分配内存 nil
	var d []int        //申明
	d = make([]int, 2) //分配内存 并且会默认填充对应类型的零值
	if c == nil {
		// 只声明一个切片 未分配内存
		t.Log("c equal nil")
	}
	if d == nil {
		// 空切片 分配内存了 make会填充默认值
		t.Log("d equal nil")
	}
	t.Logf("%v %T", c, c) //nil [] []int
	t.Logf("%v %T", d, d) //无nil [0 0] []int
	// 声明一个空切片
	var e = []int{}
	//var e []interface{}   // 声明并初始化, 只是未填充,这里用具体类型int string会报错 nil [] []interface {}
	if e == nil {
		// 空切片分配内存了
		t.Log("e equal nil")
	}
	t.Logf("%v %T", e, e) //无nil [] []int

	// 指针 创建该类型的指针 分配内存 *f 取指针指向变量的值 零值
	var f = new(int)
	if f == nil {
		t.Log("f equal nil")
	}
	t.Logf("%v %T %d", f, f, *f) // 0xc000014328 *int 0

}
