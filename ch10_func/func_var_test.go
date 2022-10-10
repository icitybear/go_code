package fn_test

import (
	"fmt"
	"testing"
)

type Data struct {
	complax  []int     // 测试切片在参数传递中的效果
	instance InnerData // 实例分配的innerData
	// 注意一下，下面的ptr的值，获取的是地址，具体地址里的值是什么，并不关心
	ptr *InnerData // 将ptr声明为InnerData的指针类型
}

// 代表各种结构体字段
type InnerData struct {
	A int
}

// 声明一个全局变量
var in = Data{
	// 测试切片在参数传递中的效果
	complax: []int{2, 3, 4},
	// 实例分配的innerData
	instance: InnerData{A: 4},
	// 将ptr声明为InnerData的指针类型
	ptr: &InnerData{A: 6},
}

//测试，结构体作为函数参数传入
func TestPassByValue(t *testing.T) {
	fmt.Printf("传入函数前, 打印结构体Data实例的值:   %+v", in)
	fmt.Printf("传入函数前, 打印结构体Data实例中instance的地址: %p", &in.instance)
	fmt.Printf("传入函数前, 打印结构体Data实例的地址: %p", &in)
	fmt.Println()
	out := passByValue(in)
	fmt.Println()
	fmt.Printf("经过函数处理后, 打印结构体Data实例的值: %+v", out)
	fmt.Printf("经过函数处理后, 打印结构体Data实例中instance的地址:   %p", &out.instance)
	fmt.Printf("经过函数处理后, 打印结构体Data实例的地址:  %p", &out)
}

//函数内部的变量，是局部变量
func passByValue(inFunc Data) Data {
	fmt.Printf("在函数内部, 打印结构体Data实例的值:   %+v", inFunc)
	fmt.Printf("在函数内部, 打印结构体Data实例中instance的地址: %p", &inFunc.instance)
	fmt.Printf("在函数内部, 打印结构体Data实例的地址: %p", &inFunc)
	return inFunc
}

type People struct {
	Name string
	Age  int8
}

func passTest(p People) {
	fmt.Printf("传递到函数中的地址: %p\n", &p)
	// 传递到函数中的地址: 0xc0000b6060
}

func TestPeople(t *testing.T) {
	myName := People{Name: "renshanwen", Age: 18}
	fmt.Printf("主函数初始化的变量地址: %p\n", &myName)
	// 主函数初始化的变量地址: 0xc0000b6048
	passTest(myName)
	fmt.Printf("不会影响到主函数的变量: %p\n", &myName)
	//不会影响到主函数的变量: 0xc0000b6048
}

func passTest2(p *People) {
	fmt.Printf("传递到函数中的指针变量的内存地址: %p，指针变量指向的内存地址：%p\n", &p, p)
	//传递到函数中的指针变量的内存地址: 0xc00000e050，指针变量指向的内存地址：0xc00000c060
	p.Name = "chengcheng" // 修改名字
}

func TestPeople2(t *testing.T) {
	myName := &People{Name: "renshanwen", Age: 18}
	fmt.Printf("主函数初始化指针变量的内存地址: %p, 指针变量指向的内存地址:%p\n", &myName, myName)
	//主函数初始化指针变量的内存地址: 0xc00000e048, 指针变量指向的内存地址:0xc00000c060
	passTest2(myName)
	fmt.Printf("被修改后的名字: %s\n", myName.Name)
	//被修改后的名字: chengcheng
	fmt.Printf("主函数被影响后指针变量的内存地址: %p, 指针变量指向的内存地址:%p\n", &myName, myName)
	//主函数被影响后指针变量的内存地址: 0xc00000e048, 指针变量指向的内存地址:0xc00000c060
	//如果返回值才会修改
}

//函数类型本质也是指针类型
func TestFunc(t *testing.T) {
	f1 := func(i int) int { return 1 }
	// c := f1(5)
	// fmt.Printf("c: %+v, 内存地址：%p", c, &c)
	fmt.Printf("f1: %T, 内存地址：%p", f1, &f1)
	f2 := f1 //赋值后 内存地址是不一样 但是f1 f2指向的值一样 +v不好打印
	fmt.Printf("f2: %T, 内存地址：%p", f2, &f2)
}
