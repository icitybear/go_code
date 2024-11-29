package my_map

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/iancoleman/orderedmap"
)

func TestInitMap(t *testing.T) {
	// var m0 map[int]string
	// m0[0] = "hello"
	// t.Log(m0)
	m1 := map[int]int{1: 1, 2: 4, 3: 9}
	t.Log(m1[2])
	t.Log(m1[10]) // 类型0值
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

var (
	VIVO_MEDIA_CODE_CLIENT_ID_MAP = map[string]map[string]map[string]string{
		"vivo": {
			"taqu": {
				"client_id":         "20210616002",
				"client_secret":     "ABE6E192C7980631881813F2B8FF1AD69294509FCF8E0680164B72AF11FE541C",
				"src_id":            "ds-202103082874",
				"refresh_redis_key": "vivo:refreshtoken",
				"token_redis_key":   "vivo:token",
			},
			"fengyue": {
				"client_id":         "20240416003",
				"client_secret":     "C8E175B188AE6190A74EFCA6C2CC937729D74A093301DEB912179132DBE43C73",
				"src_id":            "ds-202404105996", // 必须配套才不会报错
				"refresh_redis_key": "fengyue:vivo:refreshtoken",
				"token_redis_key":   "fengyue:vivo:token",
			},
		},
		"vivoWeb2": {
			"taqu": {
				"client_id":         "20230508006",
				"client_secret":     "928A9CAA71CC5E04B399C1570A0E79B04EDA53157A676F7A0845067288414557",
				"src_id":            "ds-202304064400",
				"refresh_redis_key": "taqu_vivoWeb2_refreshtoken", // vivoweb2:refreshtoken是旧规则
				"token_redis_key":   "taqu_vivoWeb2_token",        // vivoweb2:token是旧规则
			},
			"fengyue": {
				"client_id":         "20240419010",
				"client_secret":     "8C3146AA66E88A103725D33A88CE76BE130AAF3A6FFF5CB94AA960242500421F",
				"src_id":            "ds-202404194435",
				"refresh_redis_key": "fengyue_vivoWeb2_refreshtoken",
				"token_redis_key":   "fengyue_vivoWeb2_token",
			},
		},
	}
)

// 嵌套里的ok是链式的
func TestQt(t *testing.T) {
	config := VIVO_MEDIA_CODE_CLIENT_ID_MAP
	mediaCode := "vivoWeb2"
	// mediaCode = strings.ToLower(mediaCode)
	productCode := "fengyue"
	data, ok := config[mediaCode][strings.ToLower(productCode)]
	if !ok {
		fmt.Println("not ok")
		return
	}
	fmt.Printf("%s\n", mediaCode)
	fmt.Printf("ok data:%v", data)
}
