package extension

import (
	"fmt"
	"testing"
)

type Pet struct {
}

func (p *Pet) Speak() {
	fmt.Print("...")
}

func (p *Pet) SpeakTo(host string) {
	p.Speak()
	fmt.Println(" ", host)
}

type Dog struct {
	//p *Pet 指针变量 但是对应方法要重新写下 直接使用Pet类型 就不用了
	Pet //使用内嵌结构体
}

//重载 多个同名函数同时存在，具有不同的参数个数/类型。
//重写  方法与其父类有相同的名称和参数

// 不支持重载 Pet结构体 speakTo 使用的是pet指针的speak 除非重写Dog指针的 speakTo 所以说不支持重载
// 特殊其他语言 父类speakTo调用p.Speak时会重载子类dog的speak，但是go不会。
func (d *Dog) Speak() {
	fmt.Print("Wang!")
}

// 进行区分
// func (d *Dog) SpeakTo(host string) {
// 	d.Speak()
// 	fmt.Println(" ", host)
// }

// 无法支持lsp （里氏替换原则）子类交换原则 使用父类的地方都能使用子类  显示转换 （认为非继承）
func TestDog(t *testing.T) {
	dog := new(Dog)
	//不是真正的继承  内嵌了pet 理解成组合

	//内部的speak调用的还是pet的speak ...  chen
	//除非 dog重写了SpeakTo里面用了Dog的指针  Wang!  chen
	dog.SpeakTo("Chen")

	//⽗类Pet的定义的⽅法SpeakTo ⽆法访问⼦类Dog的数据和⽅法Speak
}

func TestPet(t *testing.T) {
	pet := new(Pet)
	pet.Speak()
}
