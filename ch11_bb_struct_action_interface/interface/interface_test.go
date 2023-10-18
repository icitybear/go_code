package interface_test

import (
	"fmt"
	"testing"
)

// 接口 类型断言 类型分支 接口类型断言 结构体类型断言
type Programmer interface {
	//自定义interface
	WriteHelloWorld() string
}

type GoProgrammer struct {
}

func (g *GoProgrammer) WriteHelloWorld() string {
	return "fmt.Println('Hello World')"
}

func TestClient(t *testing.T) {
	//var p Programmer  p = new(GoProgrammer)
	p := new(GoProgrammer)
	t.Log(p.WriteHelloWorld())
}

// 接口变量 Programmer是接口类型
func TestVar(t *testing.T) {
	var prog Programmer = &GoProgrammer{} //类型 数据（实例） 把实现接口的结构体的指针类型 给接口
	t.Log(prog.WriteHelloWorld())
}

// 类型断言和类型switch
func TestDy(t *testing.T) {
	var x, y interface{} // interface{} any
	x = 10.08
	value, ok := x.(float64)
	fmt.Println(value, ",", ok) //10.08,true

	x = 10.08
	value2, ok := y.(int)        // 断言失败的情况 类型推断
	fmt.Println(value2, ",", ok) //0,false

	var params map[string]interface{} // 断言的x值为nil  x 是 nil 接口值
	// panic: interface conversion: interface {} is nil, not int [recovered] 如果不传ok，要配合switch
	value3, ok := params["title"].(int)
	fmt.Println(value3, ok) // 0 false

	value4, ok := params["title"].(string)
	fmt.Println(value4, ok) // "" false
}

// 类型分支
func TestSc(t *testing.T) {
	// var a interface{}
	// a = nil
	// value3 := a.(int) // 这种写法如果类型断言时为nil会panic 除非配合switch使用
	// fmt.Println(value3)

	// var a interface{}
	// a = 10.565 // float64
	// value3 := a.(int) // 断言失败的情况 也必须要接收err panic: interface conversion: interface {} is float64, not int [recovered]
	// panic: interface conversion: interface {} is float64, not int
	// fmt.Println(value3)

	var x interface{}        // interface{} any
	x = 5656.000             //即使 为nil 也不会panic
	switch val := x.(type) { //类型断言 强转
	case int:
		fmt.Println("int is", val)
	case float64:
		fmt.Println("float64 is", val)
	case string:
		fmt.Println("string is", val)
		// 也可为类型别名 或者结构体
	default:
		fmt.Printf("Unknown type %T %+v", val, val)
	}
}

type IAnimal interface {
	GetName() string
	Call() string
	FavorFood() string
}

type Animal struct {
	Name string
}

func (a Animal) Call() string {
	return "动物的叫声..."
}

func (a Animal) FavorFood() string {
	return "爱吃的食物..."
}

func (a Animal) GetName() string {
	return a.Name
}

func NewAnimal(name string) Animal {
	return Animal{Name: name}
}

type Pet struct {
	Name string
}

func (p Pet) GetName() string {
	return p.Name
}
func NewPet(name string) Pet {
	return Pet{Name: name}
}

type Dog struct {
	animal *Animal // 继承指针类型
	pet    Pet     // 继承值类型 思考是否能继承到指针方法呢
}

func NewDog(animal *Animal, pet Pet) Dog {
	return Dog{animal: animal, pet: pet}
}

func (d Dog) FavorFood() string {
	return d.animal.FavorFood() + "骨头"
}
func (d Dog) Call() string {
	return d.animal.Call() + "汪汪汪"
}

func (d Dog) GetName() string {
	return d.pet.GetName()
}

// 结构体断言
func TestCc(t *testing.T) {
	var animal = NewAnimal("中华田园犬")
	var pet = NewPet("泰迪")
	// 必须实现接口才能赋值
	var ianimal IAnimal = NewDog(&animal, pet)
	if dog, ok := ianimal.(Dog); ok {
		fmt.Println(dog.GetName())
		fmt.Println(dog.Call())
		fmt.Println(dog.FavorFood())
	}
	// false
	if dog, ok := ianimal.(Animal); ok {
		fmt.Println(dog.GetName())
		fmt.Println(dog.Call())
		fmt.Println(dog.FavorFood())
	}

	// 父类实现了某个接口 不代表组合类它的子类也实现了这个接口
	// 这里父类实现了， 子类Dog也重写实现了 true
	if dog, ok := ianimal.(IAnimal); ok {
		fmt.Println(dog.GetName())
		fmt.Println(dog.Call())
		fmt.Println(dog.FavorFood())
	}
}
