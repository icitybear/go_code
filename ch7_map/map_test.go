package my_map

import "testing"

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
