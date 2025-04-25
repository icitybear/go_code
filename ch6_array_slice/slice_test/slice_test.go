package slice_test

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// 切片的类型字面量中只有元素的类型，没有长度
func TestArr(t *testing.T) {
	var idList []int
	idList = []int{25, 26}
	for _, id := range idList {
		t.Log(id)
	}
	// map类型
	c := map[string]string{}
	t.Logf("%T %v", c, c) // map[string]string map[]
	if c == nil {
		t.Log("hhh")
	}
	// if _, ok := c; !ok {
	// 	t.Log("ccc")
	// }

}

type stu struct {
	Name string
}

func TestSliceSm(t *testing.T) {
	var s []int
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	spew.Println(s)

	var sm []*stu // 切片可以只声明 然后直接使用 map必须声明初始化
	spew.Println(sm)
	var sm1 []*stu
	sm = append(sm, sm1...) // 空的不会追加上
	spew.Println(sm)
	fmt.Println(len(sm), cap(sm))
	sm = append(sm, &stu{Name: "csx"})
	spew.Println(sm)
}

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

// 切片自动扩容
func TestSliceGrowing(t *testing.T) {
	s := []int{}
	for i := 0; i < 10; i++ {
		s = append(s, i)
		t.Log(len(s), cap(s))
	}
}

// 数据共享 通过append扩容后解决
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

// 空切片 与 只声明
func TestSliceComparing(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{1, 2, 3, 4} //切片
	//c := [...]int{1, 2, 3, 4} 数组
	// if a == b { //切片只能和nil比较 指针与指针不能比较 计算
	// 	t.Log("equal")
	// }
	t.Log(a, b)
	var c []int        //只声明 不分配内存 nil
	var d []int        //申明
	d = make([]int, 2) //分配内存 并且会默认填充对应类型的零值
	var x = []int{}

	if c == nil {
		// 只声明一个切片 未分配内存
		t.Log("c equal nil") // 会执行
	}
	if d == nil {
		// 空切片 分配内存了 make会填充默认值
		t.Log("d equal nil") // 不执行
	}
	if x == nil {
		// 空切片 不填充默认值
		t.Log("x equal nil") // 不执行
	}
	t.Logf("%v %T", c, c) //nil [] []int
	t.Logf("%v %T", d, d) //无nil [0 0] []int
	t.Logf("%v %T", x, x) //无nil [] []int
	// 声明一个空切片
	var e = []int{}
	// 声明并初始化, 只是未填充,这里用具体类型int string会报错 nil [] []interface {}
	// var e []interface{} 语法错误， 要使用这个var e = []interface{} Compilation failed
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

// 删除
func TestDel(t *testing.T) {
	// 通过切片的切片
	slice0 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice1 := slice0[:len(slice0)-5] // 删除 slice3 尾部 5 个元素 这种不变 起始指针不变 下标0
	slice2 := slice0[5:]             // 删除 slice3 头部 5 个元素 会变化 起始指针变了 下标5
	// slice2 := append(slice0[:0], slice0[5:]...) // 2种删除头部的 返回的切片地址不一样， 这种不变 起始指针不变 下标0
	fmt.Printf("%p, %p, %p\n", slice0, slice1, slice2) // 地址是头部指针 0x1400011a050, 0x1400011a050, 0x1400011a078

	slice3 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("%p\n", slice3) // 0x1400011a0a0
	// 通过append 删除并没有自动扩容，所以不会返回新切片地址，还是旧切片
	slice4 := append(slice3[:0], slice3[3:]...) // 删除开头三个元素
	fmt.Printf("%p, %p\n", slice3, slice4)      // 0x1400011a0a0, 0x1400011a0a0
	slice5 := append(slice3[:1], slice3[4:]...) // 删除中间三个元素
	fmt.Printf("%p, %p\n", slice3, slice5)      // 0x1400011a0a0, 0x1400011a0a0
	slice6 := append(slice3[:0], slice3[:7]...) // 删除最后三个元素
	fmt.Printf("%p, %p\n", slice3, slice6)      // 0x1400011a0a0, 0x1400011a0a0
	slice7 := slice3[:copy(slice3, slice3[3:])] // 删除开头前三个元素 这种不变 起始指针不变 下标0
	fmt.Printf("%p, %p\n", slice3, slice7)      // 0x1400011a0a0, 0x1400011a0a0
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

	t.Log(s1[0])
	//v, ok := s1[0] ok  只有map才是有2个返回值（有key设定 除非s1 map[int]int） ok为对应元素是否存在

	// 还不能使用s1[0]值 必须使用s1切片 != nil 具体值 s1[0]不能跟nil比较
	// nil是一个预先声明的标识符，只能表示指针、通道、函数、接口、映射或切片
	if s1 != nil {
		t.Log(s1)
	}

	var s2 []int
	//s2 := []int{}

	//直接用现有s1
	s2 = append(s1, 5, 6)
	s2 = append(s2, 7)
	// [1, 5, 6, 7] [100, 5, 6, 7]
	t.Log(s2)
	for idx, elem := range s2 {
		t.Log(idx, elem)
	}

	s3 := s2[1:]
	t.Log(s3)
	// 改了值 s2也会变 （数据共享问题） 切片 slice或者map  切片 只是  起始地址 + 容量

	// 通过range获取数组的值 -> 不能修改原数组中结构体的值： 只是值副本 只能通过下标 科学
	s3[1] = 10 //数组会报错
	for idx, elem := range s3 {
		t.Log(idx, elem)
	}
	t.Log(s3)
	t.Log(s2) //原来s2的也变了
}

func TestEqual(t *testing.T) {
	s1 := []int{1, 3}
	t.Log(len(s1), cap(s1))

	s2 := []int{3, 1}
	t.Log(len(s2), cap(s2))
	// 切片时引用类型 不能比较
	// if s1 == s2 {
	// 	fmt.Println("equal")
	// } else {
	// 	fmt.Println("not equal")
	// }

}
