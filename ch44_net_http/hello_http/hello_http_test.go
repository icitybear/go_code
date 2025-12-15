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

func TestParse2(t *testing.T) {
	// 复杂的 URL
	// rawURL := "https://example.com/path/to/resource?name=廖康子&age=25&city=北京#section"

	// 含有特殊符号的情况 空格百分号
	rawURL := "https://admp-public.xxx.com/外拍-廖康子-荷尔蒙-口播-1216bj8啵啵-90%男生不知道-翻剪-镜像二创-xxx (5).mp4"
	// rawURL := "https://admp-public.xxx.com/外拍-廖康子-荷尔蒙-口播-1216bj8啵啵-90男生不知道-翻剪-镜像二创-xxx(5).mp4"
	encodedPath := url.PathEscape(rawURL)
	fmt.Println(0, encodedPath)
	if validator.IsUrl(rawURL) == false {
		fmt.Println("非法")
	}

	// rawURL = url.QueryEscape(rawURL) // 只限加密参数 不然会http没了
	if _, err := url.QueryUnescape(rawURL); err != nil {
		// rawURL = strings.ReplaceAll(rawURL, "%\xe7\x94", "%E7%94") // 替换为合法的 URL 编码
		// 但是可能存在很多特殊的符号，且通过链接也不能正常访问

		// URL 中的非法字符并不限于 %\xe7\x94，还包括：
		// 无效的 URL 编码（如 %g1、%1 等）。
		// 未编码的保留字符（如 ?、#、& 等）。
		// 未编码的非 ASCII 字符（如中文字符、é、ñ 等）。

		// 分割 URL 和文件名 这种不规范
		// rawURL := "https://admp-public.xxx.com/外拍-廖康子-荷尔蒙-口播-1216bj8啵啵-90%男生不知道-翻剪-镜像二创-xxx (5).mp4"
		lastSlashIndex := strings.LastIndex(rawURL, "/")
		baseURL := rawURL[:lastSlashIndex+1]
		fileName := rawURL[lastSlashIndex+1:]
		fmt.Println(baseURL)
		fmt.Println(fileName)
		rawURL = baseURL + url.QueryEscape(fileName)
		fmt.Println(1, rawURL)
	}

	// 解析 URL invalid URL escape "%\xe7\x94 特殊符号要先加密
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("解析 URL 失败:", err)
		return
	}
	domain := parsedURL.Scheme + "://" + parsedURL.Host
	fmt.Println(domain)

	fmt.Println(3, parsedURL.Path) // RequestURI() Path+RawQuery
	fmt.Println(4, parsedURL.RawPath)
	fmt.Println(5, parsedURL.RawQuery)
	fmt.Println(6, parsedURL.RequestURI())
	fmt.Println(7, parsedURL.Fragment) // section

	// 获取查询参数
	queryParams := parsedURL.Query()

	// 遍历并打印所有参数（解码后的值）
	for key, values := range queryParams {
		for _, value := range values {
			decodedValue, err := url.QueryUnescape(value)
			if err != nil {
				fmt.Println("解码失败:", err)
			} else {
				fmt.Printf("参数: %s, 解码后的值: %s\n", key, decodedValue)
			}
		}
	}
}

//------------------
// url.QueryUnescape
// url.QueryEscape 说明：
// : 被编码为 %3A
// / 被编码为 %2F
// 中文字符（如 外拍、廖康子 等）被编码为 % 开头的 UTF-8 字节序列。
// % 被编码为 %25
// 空格 被编码为 +
// ( 和 ) 被编码为 %28 和 %29
// 注意事项：
// 如果需要对整个 URL 进行编码，直接使用 url.QueryEscape 会将 :// 和 / 等 URL 结构字符也编码，这可能会导致 URL 不可用。
// 如果只想对 URL 的路径部分或参数部分进行编码，需要先拆分 URL，然后对特定部分进行编码。

//------------------
// parsedURL.Path：

// 获取 URL 的路径部分（如 /path/to/resource）。
// parsedURL.RawQuery：

// 获取 URL 的查询参数部分（如 name=廖康子&age=25）。
// parsedURL.Fragment：

// 获取 URL 的片段部分（如 section）。
// 拼接结果：

// 将路径、查询参数和片段按顺序拼接起来，形成最终的字符串。
// 注意事项：
// 如果 URL 中没有查询参数或片段，则对应的部分会被忽略。
// 如果 URL 中包含特殊字符（如 %、# 等），它们会保留在结果中。

//------------------
// 无效的 URL 编码：

// 在 URL 编码中，% 后面必须跟随两个十六进制字符（如 %20 表示空格）。
// 用户输入中的 %\xe7\x94 并不是合法的 URL 编码，因为 \xe7 和 \x94 是字节序列，而不是有效的十六进制字符。
// Go 语言的 URL 解析：

// Go 的 net/url 包在解析 URL 时会严格验证 URL 编码的合法性。
// 如果遇到无效的 % 转义序列，会返回 invalid URL escape 错误。
