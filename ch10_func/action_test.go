package fn_test

import (
	"fmt"
	"testing"
	"unsafe"
)

type Employee struct {
	Id   string
	Name string
	Age  int
}

//如果类型定义了 String() 方法，它会被用在 fmt.Printf() 中生成默认的输出：等同于使用格式化描述符 %v 产生的输出。
// 还有 fmt.Print() 和 fmt.Println() 也会自动使用 String() 方法。

func (e Employee) Act() string {
	e.Id = "001" // 不会修改值
	return fmt.Sprintf("ID:%s-Name:%s-Age:%d", e.Id, e.Name, e.Age)
}

func (e Employee) String() string {
	e.Id = "003" // 不会修改值
	return fmt.Sprintf("ID:%s-Name:%s-Age:%d", e.Id, e.Name, e.Age)
}

// 与&Employee{"0","Bob",20} 对应
// func (e *Employee) Act() string {
// 	e.Id = "001" // 指针传递会修改值
// 	return fmt.Sprintf("ID:%s/Name:%s/Age:%d", e.Id, e.Name, e.Age)
// }

// // 如果覆盖方法 String这个方法再fmt.Printf
// func (e *Employee) String() string {
// 	e.Id = "003" // 指针传递会修改值
// 	return fmt.Sprintf("ID:%s/Name:%s/Age:%d", e.Id, e.Name, e.Age)
// }

func TestStructOperations(t *testing.T) {

	e := Employee{"0", "Bob", 20}
	//e := &Employee{"0", "Bob", 20}
	fmt.Printf("Address is %v\n", e)
	t.Log(e.Act())
	fmt.Println(e)
	fmt.Printf("Address is %x\n", unsafe.Pointer(&e.Name))
	fmt.Printf("Address is %+v\n", e)
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

	e1.Name = "city"
	t.Log(e1)
	t.Log(e1.Name)
	t.Log(e1.Id) //为未填充

	t.Log(e2)
	t.Log(e2.Name)
	t.Logf("e is %T", e)
	t.Logf("e1 is %T", e1)
	t.Logf("e2 is %T", e2)
}
