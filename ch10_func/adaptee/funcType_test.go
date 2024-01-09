package adaptee_test

import (
	"fmt"
	"net/http"
	"testing"
)

// 声明和初始化 函数变量 把函数作为值保存到变量中 变量是函数类型
var fVar = func(arg int) {
	fmt.Println("fVar", arg) // 要注入的函数 自定义的实现（要注入的函数f）
}

// tag: 定义类型FT 是个（未实现的函数） 不是结构体struct 所以创建实例的时候不能用{}
// tag: 通过FT(函数变量)调用 创建FT类型实例 可以把fVar函数变量
type FT func(arg int)

// 函数名的本质就是一个指向其函数内存地址的指针常量 匿名函数

// tag: 类型实现自己的成员方法(可以去实现各种接口，对应的自定义注函数就能入这个方法了)  注意这里不使用 *FT 函数类型不能再用指针了
func (f FT) Hello(arg int) {
	// 参数 刚好是 函数类型FT需要的
	fmt.Println("FT Hello", arg)
	f(arg) // 把函数f注入了该接口IFC方法Hello的流程  f => FT需要是个函数类型
}

// FT类型又能实现接口IFC （适配器FT实现接口就行）
type IFC interface { // 参照http.Handler接口
	Hello(arg int)
}

func ClientAdaptee(ifc IFC) {
	ifc.Hello(10)
}

func TestXxx(t *testing.T) {
	// 直接调用函数类型的变量
	fVar(3) // fVar 3

	// 函数类型的自定义类 如何创建实例, 传递一个函数类型进来，FT的实例就是当成函数用
	FT(fVar)(3) // 这里的3是fVar函数的参数 FT{}是错误的  fVar 3

	// 调用FT类型的成员方法  底层一样，通过FT包装
	FT(fVar).Hello(4) // 这里的3是fVar函数的参数  FT Hello 4 与 fVar 4

	// tag: 使用了适配器FT 自定义方法fVar
	// 只依赖接口，不依赖具体实现。依赖倒置原则。ClientAdaptee只需要定义IFC接口，无需关心实现，然后自定义函数用适配器FT包装（该适配器函数实现了该接口）
	x := FT(fVar)    // var x IFC
	ClientAdaptee(x) // ClientAdaptee(FT(fVar))  这样就像Handle调用 使用了适配器HandlerFunc    // FT Hello 10 与 fVar 10
	// tag: IFC与fVar关系 用一个接口(IFC)去接收一个函数(fVar)，接口规定了函数名称FT，但是不关心函数(fVar)的实现。

	// 比如标准库 net.http包
	// type HandlerFunc func(ResponseWriter, *Request)
	// HandlerFunc类型就类似FT, 实现了http.Handler接口（http.Handle的参数要求），HandlerFunc就是适配器了，真正要执行的函数是boy（boy要实现对应接口）
	// http.Handle(pattern string, handler http.Handler)
	http.Handle("/", http.HandlerFunc(Boy))
	// HandlerFunc实现了ServeHTTP该方法 就是实现了http.Handler接口
	// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	// 	f(w, r)
	// }
}

// 实现就可以通过函数注册的方式进入到xxx的流程里面。同时还可以用于实现观察者等模式。
func Boy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("工")
	w.Write([]byte("boy"))
}
