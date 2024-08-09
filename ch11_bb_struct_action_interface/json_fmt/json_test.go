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
