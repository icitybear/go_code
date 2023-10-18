package reflect_test

import (
	"fmt"
	"reflect"
	"testing"
)

type Service struct {
	// Tag的key为json val为name
	Name string `json:"name"`
}

// Service类型  方法GetName接收者的类型是Service类型而非Service的指针类型，
// 在结构体当中，接收者类型的区别将影响该结构体方法的可见性。假设接收者类型为指针类型则该方法称为指针方法，假如是值类型，则该方法称为值方法。
func (s Service) GetName() string {
	return s.Name
}

// Service的指针类型
func (s *Service) SetName(name string) {
	// 将 ins.Name 形式转换为 (*ins).Name。方法也是一样ins.Meth() 形式转换为 (*ins).Meth()
	s.Name = name
}

func worker() {
	s1 := new(Service) //&Service{} 指针 一个指针类型拥有它以及它的基底类型为接收者类型的所有方法（包括值方法）

	fmt.Printf("s1 %T %+v\n", s1, s1)
	s2 := Service{} //实例

	s2.SetName("citybear") // 对于类型 T，如果在 *T 上存在方法 Meth()，并且 t 是这个类型的变量，那么 t.Meth() 会被自动转换为 (&t).Meth()。但是反射中是找不到这个方法的
	s2.GetName()           //基底类型却只能拥有以它本身为接收者类型的方法。
	fmt.Printf("s2 %T %+v\n", s2, s2)
	//ins是结构体指针的情况 是因为Go语言为了方便开发者访问结构体指针的成员变量，使用了语法糖（Syntactic sugar）技术，
}

//反射调用造成该恐慌的原因是rv.MethodByName()并没有返回一个reflect.Value而是一个nil  （找不到方法就是nil）
//原因在于，Service类型是*Service的基底类型，在Go的指针知识中，有一条规则：一个指针类型拥有它以及它的基底类型为接收者类型的所有方法（包括值方法），而它的基底类型却只能拥有以它本身为接收者类型的方法。（记住这个规则）

// 正常调用结构体方法
// 以接收者类型为*Service的结构体方法，其中变量s是Service的值的指针的副本，但如果接收者类型为Service的结构体方法，也就是值方法，则变量s是Service的一个副本，如果尝试改变s的副本的属性值，则对s的属性值是不会造成影响的
func worker1() {
	s := new(Service) //&Service{}
	s.SetName("citybear")
	name := s.GetName()
	fmt.Printf("normal GetName return %s\n", name)
}

// https://studygolang.com/articles/20315 反射可以将“接口类型变量”转换为“反射类型变量” ， 反射类型指的就是reflect.Type和reflect.Value
// 通常一个结构体实例 认为是Value,Type组成 Value也是实例本身
// 通过反射 一个方法根据传值动态地调用结构体的方法 经常用来 调用grpc  在Web程序框架设计中编写调度分发控制器的时候
// 反射方法传入与传出都是数组
func worker2() {
	s := Service{}
	// 获取的是实例 参数必须是指针&s
	//panic: reflect: call of reflect.Value.Call on zero Value
	rv := reflect.ValueOf(&s) //如果不调用SetName 是可以直接传s ,不然会报错，因为该方法时绑定到结构体指针类型

	// 带参数的调用
	params := []reflect.Value{reflect.ValueOf("city")} //构造一个类型为reflect.Value的切片 作为传参数组
	// 查找反射获取的实例  获取方法名xxx := refUser.MethodByName(xxx) 后续xxx.Call
	// SetName方法并调用
	rv.MethodByName("SetName").Call(params)

	// GetName方法  返回的是数组 虽然实际方法就返回一个参数 但是就是ret[0]
	// 不带参数用nil 或者 make([]reflect.Value, 0)
	ret := rv.MethodByName("GetName").Call(nil)
	//[]reflect.Value  [city]reflect  返回类型也是  类型为reflect.Value的切片
	fmt.Printf("reflect call res %T %+v \n", ret, ret)
	fmt.Printf("reflect call return %s\n", ret[0].String())
	ret2 := rv.MethodByName("GetName").Call(make([]reflect.Value, 0))
	fmt.Printf("reflect call res2 %T %+v \n", ret2, ret2)

}

// 通过反射获取结构体注释 tag
func worker3() {
	s := Service{}
	rt := reflect.TypeOf(s)
	if field, ok := rt.FieldByName("Name"); ok {
		// 获取标签的值
		tag := field.Tag.Get("json")
		fmt.Printf("field tag is %s\n", tag)
	}
}

func TestXxx(t *testing.T) {
	worker()
	worker1()
	worker2()
	worker3()
}
