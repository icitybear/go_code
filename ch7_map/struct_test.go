package my_map

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

type HookData struct {
	Name   string  `mapstructure:"name"`
	Age    int32   `mapstructure:"age"`
	Job    string  `mapstructure:",omitempty"` // 可忽略
	No     int64   `mapstructure:"no"`
	Score  float64 `mapstructure:"score"`
	Score1 string  `mapstructure:"score1"` // 源数据是数值的情况下 通过钩子转成字符串百分比
}

func TestMapToStruct(t *testing.T) {

	// jsonStr := "[{\"Name\":\"zhangsan\",\"aGe\":18},{\"Name\":\"lisi\",\"aGe\":22}]"
	jsonStr := "[{\"Name\":\"zhangsan\",\"age\":18,\"no\":10000000000001,\"score\":\"10.05678\",\"score1\":120.0567811},{\"Name\":\"lisi\",\"age\":22,\"no\":7044144249855934983,\"score\":\"20.05678\",\"score1\":230.056922}]"
	// score类型与HookData的Score 对不上的情况 希望string转float64

	// jsonStr := "{\"Name\":\"zhangsan\",\"aGe\":18}" 类型要与json.Unmarshal的对上
	// json: cannot unmarshal object into Go value of type []map[string]interface {}

	var respData []map[string]any
	// json字符串与map转换，不管是不是map切片
	// err := json.Unmarshal([]byte(jsonStr), &respData) // json int64转化丢失 7044144249855934983 变成 7044144249855935488
	// 本质 上游的json所传的id数值比较大，超过了float64的安全整数范围 （精度丢失） 解决方案1:no从int64改为string 方案2:UseNumber
	decoder1 := json.NewDecoder(strings.NewReader(jsonStr)) // 从一个流里进行解码 方式int64转换出问题
	decoder1.UseNumber()
	err := decoder1.Decode(&respData)
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Println(respData)

	// result := HookData{} // '' expected a map, got 'slice' 切片也要用切片接受 不然报错
	result := make([]*HookData, 0)
	// 包容的结构体
	var metadata mapstructure.Metadata
	// 创建一个新的解码器配置
	decoderConfig := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true, // 启用弱类型转换 WeakDecodeMetadata 本质
		Result:           &result,
		Metadata:         &metadata,
		DecodeHook:       customHook, // 设置自定义解码钩子
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 使用解码器进行解码
	err = decoder.Decode(&respData)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("metadata:%+v", metadata)
	spew.Println(result)

}

// 如果要判断字段名，就要使用个全局的，然后按字段遍历

// 自定义钩子 必要的参数 f源数据反射 t结果数据反射 data源数据
func customHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	// 只要发现有问题就会使用钩子去替换
	// f源 t目标
	// 检查源数据类型是否为 json.Number
	if from == reflect.TypeOf(json.Number("")) {
		// 检查目标类型是否为 int64
		if to.Kind() == reflect.Int64 {
			// 将 json.Number 转换为 int64
			number := data.(json.Number)
			return number.Int64() // 刚好2个值
		}
	}

	// UseNumber因为json使用了这个功能
	if from == reflect.TypeOf(json.Number("")) {
		// 检查目标类型是否为 int64
		if to.Kind() == reflect.String {
			dd, _ := data.(json.Number).Float64()
			return fmt.Sprintf("%.2f", dd*100) + "%", nil
		}
	}

	if from.Kind() == reflect.Float64 && to.Kind() == reflect.String {
		return fmt.Sprintf("%.2f", data.(float64)*100) + "%", nil
	}

	switch to.Kind() {
	// case reflect.Int64:
	// 	return strconv.ParseInt(data.(string), 10, 64)
	// 	// 如果是字符串string strconv.ParseInt(data.(string), 10, 64)
	// 实际是json.Number

	// 只要是float64的 tag:一般要2边都判断,不然类型断言panic报错
	// case reflect.Float64:
	// 	// 字段一致的情况下，类型对不上，通过判断 比如string转float64
	// 	d, _ := strconv.ParseFloat(data.(string), 64)
	// 	// 是否追加小数点
	// 	return d, nil

	// 只判断目标，是字符串类型就用csx, 源数据对应没有的字段也不会去设置
	// case reflect.String:
	// 	fmt.Println("to Float", from.Kind(), data)
	}
	return data, nil
}
