package json_fmt

import (
	"encoding/json"
	"fmt"
	"testing"
)

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
