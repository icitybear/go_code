package adaptee_test

import (
	"fmt"
	"testing"
)

// 业务场景描述：
// 通过Lightning接口连接电脑，mac实现Lightning接口，但win的电脑实现的是USB接口，此时需要通过一个适配器，lightning -> USB 如果是type-c的接口呢，增加一个适配

// 1.target interface, declare computer interface
type computer interface {
	InsertIntoLightningPort() // 当前系统业务所期待的接口，它可以是抽象类或接口。
}

// 2.implement target struct
type mac struct {
}

func (m *mac) InsertIntoLightningPort() {
	fmt.Println("Lightning connector is plugged into mac machine.")
}

// 3.adpatee 适配者类 被访问和适配的现存组件库中的组件接口。
type windows struct {
}

func (w *windows) InsertIntoUSBPort() {
	fmt.Println("USB connector is plugged into windows machine.")
}

// 4.adapter 适配器类 它是一个转换器，通过继承或引用适配者的对象，把适配者接口转换成目标接口，让客户按目标接口的格式访问适配者。
// type定义了个新类型（结构体）, 还有函数作为适配器的写法 尤其是当你只需要实现接口的一个功能时
type winAdapter struct {
	winMachine *windows
}

func (w *winAdapter) InsertIntoLightningPort() {
	fmt.Println("Adapter converts Lightning signal to USB.")
	w.winMachine.InsertIntoUSBPort()
}

// 5.client
type client struct {
}

func (c *client) insertLightningConnectorIntoComputer(com computer) {
	fmt.Println("Client inserts Lightning connector into computer.")
	com.InsertIntoLightningPort()
}

func TestX(t *testing.T) {
	client := &client{}

	mac := &mac{}

	client.insertLightningConnectorIntoComputer(mac)

	windowsMachine := &windows{}
	// tag: 适配器进行了转换 这里适配器是struct类型 也有func类型的（struct代码更繁琐点，但是方便扩展多个方法）
	windowsMachineAdapter := &winAdapter{
		winMachine: windowsMachine,
	}
	// 让客户强制使用目标接口 computer
	client.insertLightningConnectorIntoComputer(windowsMachineAdapter)

	// Client inserts Lightning connector into computer.
	// Lightning connector is plugged into mac machine.
	// Client inserts Lightning connector into computer.
	// Adapter converts Lightning signal to USB.
	// USB connector is plugged into windows machine.
}
