package my_map

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type Person struct {
	Name   string
	Age    int
	Gender string
}

type Class struct {
	No   string
	Pers Person
}

func StructToMap(obj interface{}) map[string]interface{} {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		value := objValue.Field(i).Interface()
		data[field.Name] = value
	}

	return data
}

// 原生利用反射 结构体转map
func TestStructToMap(t *testing.T) {
	person := Person{
		Name:   "张三",
		Age:    20,
		Gender: "男",
	}

	result := StructToMap(person)
	fmt.Println(result)

	cls := Class{
		No:   "1",
		Pers: person,
	}
	result2 := StructToMap(cls)
	fmt.Println(result2)
}

type Person2 struct {
	Id      int    `mapstructure:"id"`
	Name    string `mapstructure:"name"`
	Age     int    `mapstructure:"age"`
	Job     string `mapstructure:",omitempty"`
	MyClass string `mapstructure:"my_class"`
}

func TestStructToMap2(t *testing.T) {
	p := &Person2{
		Name: "dj",
		Age:  18,
	}

	var m map[string]interface{}
	mapstructure.Decode(p, &m) // 转的时候Job忽略了

	data, _ := json.Marshal(m) // map序列化成json字符串
	fmt.Println(string(data))  // {"age":18,"id":0,"my_class":"","name":"dj"}
	//tag: 考虑 有时候 并不想 struct的字段全转 又不能全忽略
	// map转struct 倒是可以动态 选择字段转
}
