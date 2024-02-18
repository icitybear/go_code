package my_map

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/iancoleman/orderedmap"
)

func TestInitMap(t *testing.T) {
	// var m0 map[int]string
	// m0[0] = "hello"
	// t.Log(m0)
	m1 := map[int]int{1: 1, 2: 4, 3: 9}
	t.Log(m1[2])
	t.Logf("len m1=%d", len(m1))
	m2 := map[int]int{}
	m2[4] = 16
	t.Logf("len m2=%d", len(m2))
	m3 := make(map[int]int, 10)
	t.Logf("len m3=%d", len(m3))
	//正确的
	m4 := make(map[int]string, 10)
	t.Logf("len m4=%d", len(m3))
	m4[0] = "hello"
	t.Log(m4) // map[0:hello]
}

func TestAccessNotExistingKey(t *testing.T) {
	// make声明且初始化分配内存
	m := make(map[int]int, 10)
	t.Log(m) // map[]
	// 声明并初始化
	m1 := map[int]string{}
	t.Logf("%T, %v", m1[1], m1[1])
	t.Log(m1[1])
	m1[2] = "0"
	t.Log(m1[2])
	m1[3] = "0"
	// 访问不存在的键值对
	if v, ok := m1[3]; ok {
		t.Logf("Key 3's value is %s", v)
	} else {
		t.Log("key 3 is not existing.")
	}
}

// 遍历无序 使用sort包 放入切片再排序
func TestTravelMap(t *testing.T) {
	m1 := map[int]int{1: 1, 2: 4, 3: 9}
	for k, v := range m1 {
		t.Log(k, v)
	}
}

// map排序无序问题
func TestOrder(t *testing.T) {
	m1 := map[string]int{
		"a":  1,
		"c":  2,
		"b":  3,
		"c1": 4,
	}
	fmt.Println(m1) // 输出的时候key自然顺序
	b, _ := json.Marshal(m1)
	fmt.Println(string(b))
	// 序列化也是一样的
	// 1 想要指定顺序 包wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	// 2 使用封装好的 *orderedmap.OrderedMap o对象实例数据结构， 塞不进listValue
	o := orderedmap.New()
	o.Set("a", 1)
	o.Set("c", 2)
	o.Set("b", 3)
	o.Set("c1", 4)
	keys := o.Keys()
	for _, k := range keys {
		v, _ := o.Get(k)
		fmt.Println(k, v)
	}
	c, _ := json.Marshal(o)
	fmt.Println(string(c))
}
