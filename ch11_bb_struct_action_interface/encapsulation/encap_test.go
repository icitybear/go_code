package encap

import (
	"fmt"
	"testing"
	"unsafe"
)

//变量名的首字母和字段首字母必须大写,首字母大写表示公开 实例,也可以叫对象 结构体关键字 struct
type Employee struct {
	Id   string
	Name string
	Age  int
	//行为
	//String() string
}

// 给结构体加行为 以下2个方法 如果是*Employee 只有new的时候&有调，如果参数没加 * 就会有

// 通常情况下为了避免内存拷贝我们使用第二种定义方式 e *Employee
// func (e *Employee) String() string {
// 	e.Name = "city"
// 	fmt.Printf("func inner Address is %x", unsafe.Pointer(&e.Name))
// 	return fmt.Sprintf("ID:%s/Name:%s/Age:%d", e.Id, e.Name, e.Age)
// }

// 在Go语言中约定成俗不实用this也不用self 给结构体加方法
func (e Employee) String() string {
	e.Name = "bear"
	fmt.Printf("func inner Address is %x\n", unsafe.Pointer(&e.Name))
	return fmt.Sprintf("ID:%s-Name:%s-Age:%d", e.Id, e.Name, e.Age)
}

func TestCreateEmployeeObj(t *testing.T) {
	e := Employee{"0", "Bob", 20}
	e1 := Employee{Name: "Mike", Age: 30}
	e2 := new(Employee) //返回指针 &Employee{}
	e2.Id = "2"
	e2.Age = 22
	e2.Name = "Rose"
	t.Log(e)
	t.Log(e.Name)
	t.Log(e1)
	t.Log(e1.Name)
	t.Log(e1.Id)

	t.Log(e2)
	t.Log(e2.Name)
	t.Logf("e is %T", e)
	t.Logf("e2 is %T", e2)
}

func TestStructOperations(t *testing.T) {
	e := Employee{"0", "Bob", 20}
	fmt.Printf("out Address is %x\n", unsafe.Pointer(&e.Name))
	t.Log(e.String()) //返回的是方法返回的字符串
	t.Log(e.Name)
}
