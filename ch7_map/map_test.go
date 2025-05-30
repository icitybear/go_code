package my_map

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/davecgh/go-spew/spew"
	"github.com/duke-git/lancet/v2/structs"
	"github.com/iancoleman/orderedmap"
)

func TestInitMap(t *testing.T) {
	// var m0 map[int]string
	// m0[0] = "hello"
	// t.Log(m0)
	m1 := map[int]int{1: 1, 2: 4, 3: 9}
	t.Log(m1[2])
	t.Log(m1[10])                // 类型0值
	t.Logf("len m1=%d", len(m1)) // 3
	m2 := map[int]int{}
	t.Logf("len m2=%d", len(m2)) // 0
	m2[4] = 16
	t.Logf("len m2=%d", len(m2)) // 1
	m3 := make(map[int]int, 10)
	t.Logf("len m3=%d", len(m3)) // 0
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

	// 2 使用封装好的 *orderedmap.OrderedMap o对象实例数据结构， 塞不进listValue （不支持对象）
	o := orderedmap.New()
	o.Set("a", 1)
	o.Set("c", 2) // 按set顺序排序
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

type MyStruct struct {
	// 定义结构体字段
	Name        string
	Age         int
	AvgScore    string `json:"avgScore"`
	SumScore    string `json:"SumScore"`
	ClassName   string `json:"class_name"`
	CreatedTime int
}

func UnderscoreToCamel(s string) string {
	var builder strings.Builder
	upperNext := false
	firstLetter := true
	for _, r := range s {
		if r == '_' {
			upperNext = true
		} else {
			if upperNext {
				builder.WriteRune(unicode.ToUpper(r))
				upperNext = false
			} else {
				if firstLetter {
					builder.WriteRune(unicode.ToUpper(r))
					firstLetter = false
				} else {
					builder.WriteRune(r)
				}
			}
		}
	}

	return builder.String()
}
func TestOrderStruct(t *testing.T) {
	result := []MyStruct{
		// 添加结构体实例
		{"name1", 1, "10", "20", "class1", 1744775442},
		{"name2", 2, "10", "20", "class2", 1744775442},
		{"name3", 3, "10", "20", "class3", 1744775442},
	}
	// tag: 最好与结构体字段名一致 不是以json为准 （最多转下大驼峰命名 匹配字段）
	customSortSlice := []string{"name", "age", "avgScore", "SumScore", "class_name", "createdTime"}
	customTitleMap := map[string]string{
		"name":        "名称",
		"age":         "年龄",
		"avgScore":    "平均",
		"SumScore":    "总分",
		"class_name":  "课程名",
		"createdTime": "创建时间",
	}

	tempMap := make([]*orderedmap.OrderedMap, 0)
	for _, item := range result {
		// 使用orderedmap 不然得定义个排序结构体，方便序列化时有序
		o := orderedmap.New()          // o.Set 设置字段标题 跟值 支持map排序的
		targetObj := structs.New(item) // 提供操作 struct, tag, field 的相关函数
		// fmt.Printf("%+v", targetObj)   // 转为对象都是对应的字段名
		// &{
		// raw:{Name:name1 Age:1 AvgScore:10 SumScore:20 ClassName:class1 CreatedTime:1744775442}
		// rtype:0x100b67040
		// rvalue:{typ:0x100b67040 ptr:0x140001009b0 flag:153}
		// TagName:json
		// }

		// 显示个性标题
		for _, dims := range customSortSlice {
			// dims是自己想显示的 key 对应有map设置， 然后映射到结构体对应的字段
			fieldDims := UnderscoreToCamel(dims)
			dimsName := customTitleMap[dims]
			fieldDimsVal, ok := targetObj.Field(fieldDims)
			if ok {
				dimsValue := fieldDimsVal.Value()
				if fieldDims == "CreatedTime" {
					// tag: 类似元素为数组的情况，时间戳转时间字符串的情况
					if createdTime, ok := fieldDimsVal.Value().(int); ok {
						dimsValue = time.Unix(int64(createdTime), 0).In(time.Local).Format("2006-01-02 15:04:05")
					}
				}
				o.Set(dimsName, dimsValue)
			} else {
				o.Set(dimsName, "")
			}
		}

		// for _, targetKey := range sortSlice {
		// 	targetKey2 := util.UnderscoreToCamel(targetKey) // 下划线转驼峰因为结果集是结构体
		// 	targetName := titleMap[targetKey]
		// 	if mField, ok := targetObj.Field(targetKey2); ok {
		// 		o.Set(targetName, mField.Value())
		// 	}
		// }

		tempMap = append(tempMap, o)
	}

	spew.Println(tempMap)
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
