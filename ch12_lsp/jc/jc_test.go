package jc

import (
	"fmt"
	"testing"
)

type Animal struct {
	Name string
}

func (a *Animal) Call() string {
	return "动物的叫声..."
}

// 指针方法
func (a *Animal) FavorFood() string {
	return "爱吃的食物..."
}

func (a Animal) GetName() string {
	return a.Name
}

// 定义 继承自该类型的子类 Dog
type Dog struct {
	Animal // 值类型 或者 指针类型 还可以起别名
	// 继承值类型是否也能继承指针方法
}

func TestJc(t *testing.T) {
	animal := Animal{"中华田园犬"}
	// 可以在 Dog 实例上访问所有 Animal 类型包含的属性和方法
	dog := Dog{animal}

	fmt.Println(dog.GetName())          // 中华田园犬
	fmt.Println(dog.Animal.Call())      // 动物的叫声...
	fmt.Println(dog.Animal.FavorFood()) // 爱吃的食物...
}

type Pet struct {
	Name string
}

func (p Pet) GetName() string {
	return p.Name
}

type Dog2 struct {
	*Animal // 继承指针类型
	Pet     // 继承值类型 思考是否能继承到指针方法呢
}

func TestJc2(t *testing.T) {
	animal := Animal{"中华田园犬"}
	pet := Pet{"宠物狗"}
	// 要传入指针引用
	dog := Dog2{&animal, pet}

	// fmt.Println(dog.GetName())        // 提示错误ambiguous selector，方法名重复了，需要指定
	fmt.Println(dog.Animal.GetName())   // 中华田园犬
	fmt.Println(dog.Animal.Call())      // 动物的叫声...
	fmt.Println(dog.Call())             // 动物的叫声... 因为只有Animal有定义该方法，不会有选择问题
	fmt.Println(dog.Animal.FavorFood()) // 爱吃的食物...
	fmt.Println(dog.FavorFood())        // 爱吃的食物...
}
