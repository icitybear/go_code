package my_map

import (
	"fmt"
	"reflect"
	"testing"
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
