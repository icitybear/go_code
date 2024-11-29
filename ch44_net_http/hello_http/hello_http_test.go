package hello_http

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/validator"
)

// func main() {

// 	//内置http服务
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello World!")
// 	})
// 	http.HandleFunc("/time/", func(w http.ResponseWriter, r *http.Request) {
// 		t := time.Now()
// 		timeStr := fmt.Sprintf("{\"time\": \"%s\"}", t)
// 		w.Write([]byte(timeStr))
// 	})
// 	//http://172.17.10.178:10110/  因为使用的是wsl2
// 	http.ListenAndServe(":10110", nil) //监听的端口 Handler
// }

// HttpBuildQuery 构建url
func HttpBuildQuery(params map[string]interface{}) string {
	values := url.Values{}
	for key, value := range params {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		case int:
			values.Add(key, fmt.Sprintf("%d", v))
		case []string:
			for i, s := range v {
				values.Add(fmt.Sprintf("%s[%d]", key, i), s)
			}
		case []map[string]interface{}:
			for i, m := range v {
				for subKey, subValue := range m {
					subKey = fmt.Sprintf("%s[%d][%s]", key, i, subKey)
					switch sv := subValue.(type) {
					case string:
						values.Add(subKey, sv)
					case int:
						values.Add(subKey, fmt.Sprintf("%d", sv))
					case []string:
						for j, ss := range sv {
							values.Add(fmt.Sprintf("%s[%d]", subKey, j), ss)
						}
					default:
						fmt.Println("未处理的类型:", reflect.TypeOf(sv))
					}
				}
			}
		default:
			fmt.Println("未处理的类型:", reflect.TypeOf(v))
		}
	}
	return values.Encode()
}

func convertValueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case []string:
		return fmt.Sprintf(`["%s"]`, strings.Join(v, `","`))
	case []map[string]interface{}:
		var parts []string
		for _, m := range v {
			tmp, _ := json.Marshal(m)
			parts = append(parts, string(tmp))
		}
		return fmt.Sprintf(`[%s]`, strings.Join(parts, `","`))
	default:
		return ""
	}
}

func convertMapToString(data map[string]interface{}) string {
	var parts []string
	for key, value := range data {
		parts = append(parts, fmt.Sprintf("%s=%s", key, convertValueToString(value)))
	}
	return strings.Join(parts, "&")
}
func TestMain(t *testing.T) {

	postData := map[string]interface{}{
		"advertiser_id": 111,
		"data_topic":    "MATERIAL_DATA",
		"page":          1,
		"page_size":     100,
		"dimensions":    []string{"stat_time_day", "cdp_promotion_id", "material_id"},
		"metrics": []string{
			"active_register",
			"active_rate",
			"active_cost",
			// 添加其他指标
		},
		"filters": []map[string]interface{}{
			{
				"field":    "image_mode",
				"type":     1,
				"operator": 7,
				"values":   []string{"2", "3", "5", "15", "16"},
			},
			// 添加其他过滤条件
		},
		"start_time": "2024-01-01",
		"end_time":   "2024-01-01",
		"order_by": []map[string]interface{}{
			{
				"field": "stat_cost",
				"type":  "DESC",
			},
			// 添加其他排序条件
		},
	}

	postDataJson, _ := json.Marshal(postData)
	fmt.Println("序列化后的数据:", string(postDataJson))
	//

	queryString := HttpBuildQuery(postData)
	fmt.Println("生成的查询字符串:", queryString)

	queryString2 := convertMapToString(postData)
	fmt.Println("生成的查询字符串2:", queryString2)

}

func TestParse(t *testing.T) {
	// str := "素材所属主体与开发者主体不一致无法获取URL"
	str := "https://www.baidu.com"
	res, err := url.Parse(str)
	if err != nil {
		fmt.Printf("err:%+v", err)
	}
	fmt.Println(res)

	// str2 := "https://https://https://../google.com"
	// str2 := "https://www.baidu.com"
	str2 := "www.baidu.com"
	// str2 := "素材所属主体与开发者主体不一致无法获取URL" // false
	fmt.Println(validator.IsUrl(str2))
	// validURL := govalidator.IsURL(str2)
	// fmt.Printf("%s is a valid URL : %v \n", str, validURL)
}
