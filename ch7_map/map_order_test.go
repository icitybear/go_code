package my_map

import (
	"encoding/json"
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

// map转成*structpb.Value 类型， 排序是key自然
func TestOrder2(t *testing.T) {
	// 定义结构体切片
	// structs := []MyStruct{
	// 	// 添加结构体实例
	// 	{"a", 1},
	// 	{"b", 2},
	// }

	// 使用结构体切片是不可行的 如果map是希望有排序的 使用orderedmap包就变成该结构体类型了
	structs := make([]map[string]interface{}, 0)
	structs = append(structs, map[string]interface{}{"Name": "a", "Age": 1})
	structs = append(structs, map[string]interface{}{"Name": "b", "Age": 2})

	// 转换为 []*structpb.Value 类型
	values := make([]*structpb.Value, len(structs))
	for i, s := range structs {
		// 将结构体转换为 *structpb.Value 类型
		// 注意问题：structpb.NewValue函数无法直接将结构体转换为*structpb.Value类型。
		// structpb.NewValue函数只能接受基本类型（如字符串、整数、布尔等）作为参数。
		value, err := structpb.NewValue(s) // 不支持, 改成map就支持了
		if err != nil {
			fmt.Println(i, err.Error()) //错误proto: invalid type: my_map.MyStruct
		}
		values[i] = value
	}

	// 将 []*structpb.Value 类型切片转换为 []interface{} 类型切片
	interfaceSlice := make([]interface{}, len(values))
	for i, v := range values {
		fmt.Println(i, v)
		interfaceSlice[i] = v // v是*structpb.Value
	}

	// 使用 structpb.NewList 将 []interface{} 类型切片包装为 *structpb.List 类型
	list, err := structpb.NewList(interfaceSlice) // 这里的值就不会是interface了
	if err != nil {
		fmt.Println("newList", err.Error()) //newList proto: invalid type: *structpb.Value
	}
	fmt.Println(list)
}

// 如果是listValue是proto框架支持的序列化类型，但是orderedmap结构体是不支持的
// 可执行 // 如果您想将结构体转换为*structpb.Value类型，您需要手动创建一个*structpb.Value对象，并为其设置相应的字段。下面是一个示例代码：
func TestOrder3(t *testing.T) {
	// 定义结构体实例
	s := map[string]interface{}{
		"Name": "a",
		"Age":  1,
		"Bool": true,
	}

	// 创建一个新的 *structpb.Value 对象
	value := &structpb.Value{
		Kind: &structpb.Value_StructValue{
			StructValue: &structpb.Struct{},
		},
	}

	// 将结构体实例转换为 *structpb.Value_StructValue
	structValue, err := structpb.NewStruct(s)
	if err != nil {
		// 处理错误
		fmt.Println("NewStruct", err.Error())
	}

	// 将 *structpb.Struct 设置为 *structpb.Value 的字段
	value.GetStructValue().Fields = structValue.Fields

	// 使用 value 进行后续操作
	// ...

	fmt.Println(value)
	c, _ := json.Marshal(value) // data, _ = json.MarshalIndent(personMap, "", "  ")
	fmt.Println(string(c))
}
