package json_fmt

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// 结构体序列化json 反射里tag标签
type Student struct {
	Name string
	Age  int
	Ar   Address
	id   int
}

// 只有结构体内变量必须首字母大写，才可被导出的字段转化输出
type Address struct {
	Road     string `json:"road"` //反引号注释 json相关的自动时 换成小写的或者别名
	Street   string
	City     string
	Province string
	Country  string
}

func TestJsonStruct(t *testing.T) {
	// 结构体初始化 实例化 &Student 一样能json序列化
	Zhang3 := Student{
		Name: "张三",
		Age:  18,
		Ar: Address{
			Road:     "south road",
			Street:   "123 street",
			City:     "cs",
			Province: "hn",
			Country:  "CN",
		},
		id: 100, //只有结构体内变量必须首字母大写，才可被导出的字段转化输出
	}
	fmt.Printf("%T  \n", Zhang3)
	//Marshal() 函数返回编码后的 JSON （字节切片）和 error 信息（如果出错的话）。然后我们可以打印 JSON 字符串。
	//Info_of_Zhang3, err := json.MarshalIndent(Zhang3, "", "    ")
	//函数解析：第二个参数指定每行输出的开头的字符串。输出的开头，第三个参数指定每行要缩进的字符串。此时

	Info_of_Zhang3, err := json.Marshal(Zhang3)
	fmt.Printf("%T  \n", Info_of_Zhang3) //看返回时是字节切片 []uint8  []byte
	if err == nil {
		fmt.Println(string(Info_of_Zhang3))
	} else {
		fmt.Println(err)
		return
	}

	// 因为json.UnMarshal() 函数接收的参数是字节切片
	// 所以需要把JSON字符串转换成字节切片。jsonData := []byte(json_str)
	var actress Student
	err2 := json.Unmarshal(Info_of_Zhang3, &actress)
	if err2 != nil {
		fmt.Println("error:", err2)
		return
	}
	fmt.Printf("姓名：%s\n", actress.Name)
	fmt.Printf("姓名：%s\n", actress.Ar.Road)
	// 如果是切片就range打印
}

// 接口转为 JSON 格式
type Student2 map[string]interface{}
type Address2 map[string]interface{}

func TestJsonInterface(t *testing.T) {
	// 结构体能提前注释 小写字段别名字段 只有结构体内变量必须首字母大写，才可被导出的字段转化输出(map不需要)
	Zhang3 := Student2{
		"Name": "张三",
		"Age":  18,
		"Ar": Address2{
			"Road":     "renmin south road",
			"Street":   "123 street",
			"City":     "cs",
			"Province": "hn",
			"Country":  "CN",
		},
		"Year":       2022,
		"GraduateAt": 2026,
		"ceshi":      100,
	}

	InfoOfZhang3, err := json.MarshalIndent(Zhang3, "", "    ")
	if err == nil {
		fmt.Println(string(InfoOfZhang3))
	} else {
		fmt.Println(err)
	}

}

// json字符串 解析int64 float64精度丢失
func TestDecoder(t *testing.T) {
	var request = `{"id":7044144249855934983}`

	var test1 interface{}

	_ = json.Unmarshal([]byte(request), &test1) // 解码

	objStr1, err := json.Marshal(test1) // 进行编码
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(objStr1)) // {"id":7044144249855934983}
	obj := test1.(map[string]interface{})
	id := obj["id"]
	fmt.Printf("%T, %#v \n", id, id) // float64, 7.044144249855935e+18

	// 换成NewDecoder
	var test interface{}
	decoder := json.NewDecoder(strings.NewReader(request)) // 从一个流里进行解码
	decoder.UseNumber()
	err = decoder.Decode(&test) // 解码
	if err != nil {
		fmt.Println("error:", err)
	}

	objStr, err := json.Marshal(test) // 进行编码
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(string(objStr)) // {"id":7044144249855934983}
}

// josn标签里的大小写 不影响json字符串里字段的解析（这里大小写不影响）
type Student1 struct {
	Name string `json:"name"`
	Age  int    `json:"AgE"`
}

func TestXxx(t *testing.T) {
	stu := "{\"Name\":\"zhangsan\",\"aGe\":18}"
	var s Student1
	err := json.Unmarshal([]byte(stu), &s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s.Name, s.Age)
}

// json字符串数组
func TestXxx2(t *testing.T) {
	stu := "[{\"Name\":\"zhangsan\",\"aGe\":18},{\"Name\":\"lisi\",\"aGe\":22}]"
	var s []Student1
	err := json.Unmarshal([]byte(stu), &s)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range s {
		fmt.Println(v.Name, v.Age)
	}

	var unMap []map[string]interface{}
	err = json.Unmarshal([]byte(stu), &unMap)
	fmt.Println(unMap)
}

// 标签重名的影响 都不处理该标签字段
type StudentX struct {
	Name string `json:"name"`
	Age  int    `json:"AgE"`
	Tech string `json:"name"`
}

func TestR(t *testing.T) {
	stu := "{\"Name\":\"zhangsan\",\"aGe\":18}"
	var s StudentX
	err := json.Unmarshal([]byte(stu), &s)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Logf("%+v", s) // {Name: Age:18 Tech:} 多个重名的name标签导致不知道赋值给哪个都不赋值 {Name:zhangsan Age:18 Tech:}

	d := &StudentX{
		Name: "cccc",
		Age:  18,
		Tech: "ttt",
	}
	jStr, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jStr)) // {"AgE":18} 多个重名的name标签导致不知道序列化哪个 都不处理 {name":"cccc","AgE":18,"name1":"ttt"}
}

type Person struct {
	ID   uint
	Name string
	address
}

type Person2 struct {
	ID      uint
	Name    string
	Address address
}

type address struct {
	Code   int
	Street string
}

func TestMarshalPerson(t *testing.T) {
	p := Person{
		ID:   1,
		Name: "Bruce",
		address: address{ // 小写的成员，地址结构
			Code:   100,
			Street: "Main St",
		},
	}
	p2 := Person2{
		ID:   1,
		Name: "Bruce",
		Address: address{ // 小写的成员，地址结构
			Code:   100,
			Street: "Main St",
		},
	}

	fmt.Println(p.Code, p.Street)                   // 1. 组合结构体时 直接通过 Person 访问地址成员时的感觉，即地址成员似乎直接成为了 Person 的成员
	fmt.Println(p2.Address.Code, p2.Address.Street) // 100 Main St
	output, _ := json.MarshalIndent(p, "", "  ")    // 美化输出
	println(string(output))
	output2, _ := json.MarshalIndent(p2, "", "  ") // 输出没有扁平化
	println(string(output2))
}

// 序列化的结果也扁平化
// {
//   "ID": 1,
//   "Name": "Bruce",
//   "Code": 100,
//   "Street": "Main St"
// }
// 输出没有扁平化
// {
//   "ID": 1,
//   "Name": "Bruce",
//   "Address": {
//     "Code": 100,
//     "Street": "Main St"
//   }
// }

func TestUnmarshalPerson(t *testing.T) {
	str := `{"ID":1,"Name":"Bruce","address":{"Code":100,"Street":"Main St"}}`
	// 本质 嵌入式的结构体里成员就相当于提到最外层了
	var p Person
	_ = json.Unmarshal([]byte(str), &p)                             // 3. 反序列化时，私有组合对象的公共成员又不被解析
	fmt.Printf("%+v\n", p)                                          // {ID:1 Name:Bruce address:{Code:0 Street:}}
	str1 := `{"ID":1,"Name":"Bruce","Code":100,"Street":"Main St"}` // 4. 扁平化的反而能识别到成员变量 反序列化时，私有组合对象的公共成员要这样才能被解析
	_ = json.Unmarshal([]byte(str1), &p)
	fmt.Printf("%+v\n", p) // {ID:1 Name:Bruce address:{Code:100 Street:Main St}}

	// 5. "address" 与 "Address" 无区别 因为不区分大小写
	strNew := `{"ID":1,"Name":"Bruce","address":{"Code":100,"Street":"Main St"}}`
	// strNew := `{"ID":1,"Name":"Bruce","Address":{"Code":100,"Street":"Main St"}}`
	// strNew := `{"ID":1,"Name":"Bruce","Code":100,"Street":"Main St"}` // 6. 反序列化结果 ID:1 Name:Bruce Address:{Code:0 Street:}
	var p2 Person2
	_ = json.Unmarshal([]byte(strNew), &p2)
	fmt.Printf("%+v\n", p2) // {ID:1 Name:Bruce Address:{Code:100 Street:Main St}}
}

type Person3 struct {
	ID   int
	Name string
	Age  *int
}

func TestFoo(t *testing.T) {

	p := Person3{
		ID:   1,
		Name: "Bruce",
		Age:  new(int),
	}
	*p.Age = 20
	fmt.Println(p)
	// 如果p := &Person{} ,这里就是p.Age = 20

	p2 := &Person3{
		ID:   1,
		Name: "HHH",
	}
	// &{1 HHH <nil>}

	p2.Age = new(int) // &{1 HHH 0x14000020230}
	fmt.Println(p2)
}
