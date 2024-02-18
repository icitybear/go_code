package reflect_test

import (
	"fmt"
	"reflect"
	"testing"
)

func typeValue(f interface{}) {
	fmt.Println(reflect.TypeOf(f), reflect.ValueOf(f))
	fmt.Println(reflect.ValueOf(f).Type())
	fmt.Println(reflect.TypeOf(f).Kind())  //reflect.Ptr
	fmt.Println(reflect.ValueOf(f).Kind()) //reflect.Struct
}

type User struct {
	Name string
	Age  int
}

func (u User) GetName() string {
	return "ccc"
}

func (u *User) SetName(name string) {
	u.Name = name
}

func TestTypeAndValue(t *testing.T) {
	//var f int64 = 10
	s := &User{
		Name: "city",
		Age:  10,
	}
	// 	struct { Name string; Age int } {city 10}
	// struct { Name string; Age int }
	// struct
	// struct

	// *struct { Name string; Age int } &{city 10}
	// *struct { Name string; Age int }
	// ptr
	// ptr

	// 	*reflect_test.User &{city 10}
	// *reflect_test.User
	// ptr
	// ptr
	typeValue(s)

	// meth := reflect.ValueOf(s).MethodByName("SetName") //&s 与 s的区别
	// params := []reflect.Value{reflect.ValueOf("city")}
	// ret := meth.Call(params)

	// 私有方法也是不能调用的
	meth := reflect.ValueOf(s).MethodByName("GetName") //都能调用
	// make([]reflect.Value, 0)
	ret := meth.Call(nil)
	fmt.Printf("ret %T  %+v \n", ret, ret)
}

func CheckType(v interface{}) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Float32, reflect.Float64:
		fmt.Println("Float")
	case reflect.Int, reflect.Int32, reflect.Int64:
		fmt.Println("Integer")
	default:
		fmt.Println("Unknown", t)
	}
}

func TestBasicType(t *testing.T) {
	var f float64 = 12
	CheckType(f)
	CheckType(&f)
}

type Employee struct {
	EmployeeID string
	Name       string `format:"normal"`
	Age        int
}

func (e *Employee) UpdateAge(newVal int) {
	e.Age = newVal
}

func TestInvokeByName(t *testing.T) {
	e := &Employee{"1", "citybear", 30} //实例
	//按名字获取成员 reflect.ValueOf返回一个 FieldByName返回2个

	//reflect.Value的FieldByName
	t.Logf("Name: value(%[1]v), Type(%[1]T) ", reflect.ValueOf(*e).FieldByName("Name")) //Name: value(Mike), Type(reflect.Value)

	//reflect.Type的FieldByName
	if nameField, ok := reflect.TypeOf(*e).FieldByName("Name"); !ok {
		t.Error("Failed to get 'Name' field.")
	} else {
		//获取成员field 标记
		t.Log(nameField) //{Name  string format:"normal" 16 [1] false}
		t.Log("Tag:format", nameField.Tag.Get("format"))
	}

	reflect.ValueOf(e).MethodByName("UpdateAge").
		Call([]reflect.Value{reflect.ValueOf(1)}) //UpdateAge(1)
	t.Log("Updated Age:", e)
}

type Customer struct {
	CookieID string
	Name     string
	Age      int
}

func TestDeepEqual(t *testing.T) {

	a := map[int]string{1: "one", 2: "two", 3: "three"}
	b := map[int]string{1: "one", 2: "two", 3: "three"}
	//	t.Log(a == b) //map和切片 不能直接比较==  只能更nil比较 a == nil
	t.Log("a==b?", reflect.DeepEqual(a, b))

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{2, 3, 1}

	t.Log("s1 == s2?", reflect.DeepEqual(s1, s2))
	t.Log("s1 == s3?", reflect.DeepEqual(s1, s3))

	c1 := Customer{"1", "Mike", 40}
	c2 := Customer{"1", "Mike", 40}
	fmt.Println(c1 == c2)
	fmt.Println(reflect.DeepEqual(c1, c2))
}

type Vehicle struct {
	ID          int    `json:"id"`
	CityName    string `json:"city_name"`
	Provider    string `json:"provider"`
	PlateNumber string `json:"plate_number"`
	MaxPeople   int64  `json:"max_people"`
	PowerType   string `json:"power_type"`
}

func TestStructTag(t *testing.T) {
	// 使用reflect.StructTag解析这段文本的tag内容
	tag := reflect.StructTag(`json:"foo,omitempty" xml:"foo"`)
	// 直接使用Get获取json定义
	value := tag.Get("json")
	fmt.Printf("value: %q\n", value)

	reflectType := reflect.ValueOf(Vehicle{}).Type()
	fmt.Printf("fields number: %v\n", reflectType.NumField())

	for i := 0; i < reflectType.NumField(); i++ {
		fmt.Printf("%v", reflectType.Field(i).Name)
		fmt.Printf("  tag:%v\n", reflectType.Field(i).Tag.Get("json"))
	}
}
