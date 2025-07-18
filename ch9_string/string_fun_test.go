package string_test

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/Knetic/govaluate"
)

// 拼接
func TestJoinFn(t *testing.T) {
	s := "A,B,C"
	// s := "" // tag: 注意点
	parts := strings.Split(s, ",")
	if parts == nil {
		t.Log("parts nil") // 不为nil
	}
	if len(parts) == 0 {
		t.Log("parts 0") // 不为0
	}
	fmt.Println(parts, len(parts)) // 结果是[] 1  0空字符串
	// 所以会循环一次
	for key, part := range parts {
		fmt.Println(key)
		t.Log(part)
	}
	t.Log(strings.Join(parts, "-"))
}

// 分割
func TestSplitFn(t *testing.T) {
	s1 := "chihuo@golang"
	arr := strings.Split(s1, "@")
	fmt.Printf("arr is %v\n", arr)
}

// 去掉2边空格
func TestTrimFn(t *testing.T) {
	s1 := " chihuo@golang \n"
	s2 := strings.TrimSpace(s1)
	fmt.Printf("trim space '%s'\n", s2)
}

// 是否含有前缀 后缀
func TestHas(t *testing.T) {
	s1 := "chihuo@golang"
	if strings.HasPrefix(s1, "chihuo") {
		fmt.Printf("%s has prefix chihuo\n", s1)
	}
	if strings.HasPrefix(s1, "@") {
		fmt.Printf("%s has prefix chihuo\n", s1)
	}
	if strings.HasSuffix(s1, "golang") {
		fmt.Printf("%s has suffix golang\n", s1)
	}

	// str := "外拍-沈铭炜-剧情-四月底wh2桃子0501-男用户-穿的很特别-翻剪-镜像二创-他趣"
	str := "web.business.image/202410085d0dbb6ef08fd6a64ab986e"
	fmt.Println(len(str))
	fmt.Println((utf8.RuneCountInString(str)))
}

// 子字符串 直接通过切片下表获取 因为字符船 是不可改变的值类型 又是 []types
func TestSubstr(t *testing.T) {
	s1 := "chihuo@golang"
	s2 := s1[6:len(s1)]
	fmt.Printf("sub string is %s\n", s2)
}

// 转化
func TestConv(t *testing.T) {
	s := strconv.Itoa(10)
	t.Log("str" + s)
	if i, err := strconv.Atoi("10"); err == nil {
		t.Log(10 + i)
	} else {
		t.Log(i)
	}
}

// + 号
func j1() {
	s1 := "chihuo"
	s2 := "golang"
	s3 := s1 + "@" + s2
	fmt.Printf("s1 + s2 = %s\n", s3)
}

// fmt格式化Sprintf
func j2() {
	s1 := "chihuo"
	s2 := "golang"
	s3 := fmt.Sprintf("%s@%s", s1, s2)
	fmt.Printf("s1 + s2 = %s\n", s3)
}

// strings.Join
func j3() {
	s0 := strings.Join([]string{}, ",")
	fmt.Println(len(s0), s0)

	s1 := "chihuo"
	s2 := "golang"
	s3 := strings.Join([]string{s1, s2}, "@")
	fmt.Printf("s1 + s2 = %s len:%d\n", s3, len(s3))

	var str []string
	s4 := strings.Join(str, "@")
	fmt.Printf("s4 = %s len:%d\n", s4, len(s4))

	str1 := []string{}
	s5 := strings.Join(str1, "@")
	fmt.Printf("s5 = %s len:%d\n", s5, len(s5))
}

// bytes.Buffer
func j4() {
	var bt bytes.Buffer
	s1 := "chihuo"
	s2 := "golang"
	bt.WriteString(s1)
	bt.WriteString("@")
	bt.WriteString(s2)

	s3 := bt.String() //把Buffer缓存里的转成字符串
	fmt.Printf("s1 + s2 = %s\n", s3)
}

// strings.Builder
func j5() {
	var builder strings.Builder
	s1 := "chihuo"
	s2 := "golang"
	builder.WriteString(s1)
	builder.WriteString("@")
	builder.WriteString(s2)
	s3 := builder.String()
	fmt.Printf("s1 + s2 = %s\n", s3)
	builder.WriteString("citybear apend")
	s4 := builder.String()
	fmt.Println(s4)
}

// 字符串拼接
func TestXxx(t *testing.T) {
	// 一次性执行完毕的流程 用+ fmt都无所谓，循环脚本，后台挂起的还是用效率高的
	// j1()
	// j2()
	j3()
	// 5 > 4 >321
	// j5()
}

// 字字符串查询
func Stripos(haystack string, needle string, offset ...int) int {
	off := 0
	if len(offset) > 0 {
		off = offset[0]
	}
	if off > len(haystack) || off < 0 {
		return -1
	}
	// 全转为小写
	haystack = strings.ToLower(haystack[off:])
	needle = strings.ToLower(needle)
	index := strings.Index(haystack, needle) // strings.Contains(info.Content, actionName)
	if index != -1 {
		return off + index
	}
	return index
}

type Value struct {
	Name  string
	Value int32
}

func TestFmt(t *testing.T) {
	v1 := Value{
		Name:  "val",
		Value: 10,
	}
	s1 := fmt.Sprintf("%d %v %+v %#v %T %p %f", v1.Value, v1, v1, v1, v1, &v1, float64(v1.Value))
	fmt.Printf("format is %s\n", s1)
}

func TestSplit(t *testing.T) {
	// str := "1,"
	str := ""
	tmp := strings.Split(str, ",") // tag:就算返回空数组 长度也是1
	fmt.Printf("tmp is:%+v str_len:%d arr_len:%d \n", tmp, len(str), len(tmp))
	arr := []string{}
	// var arr = make([]string, 1)
	fmt.Printf("arr is:%+v  arr_len:%d \n", arr, len(arr))

	arr1 := []string{}
	// arr1 := []string{"1", "v2", "", "4"}
	jnStr := strings.Join(arr1, ",")
	fmt.Printf("jnStr is:%+v  jnStr_len:%d \n", jnStr, len(jnStr))
}

func TestJoin(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	str := ","
	var strArr []string
	for _, v := range arr {
		strArr = append(strArr, fmt.Sprintf("%d", v))
		// strArr = append(strArr, fmt.Sprintf("'%d'", v)) 特殊写法
	}
	res := strings.Join(strArr, str)
	fmt.Println(res)

	strArr1 := []string{"taqu", "miyou", "4"}
	res1 := strings.Join(strArr1, str)
	fmt.Println(res1)
}

// 在Go语言中，字符串本身并不直接支持修改其内部字符（因为字符串在Go中是不可变的），包括将字符串的首字母小写
func TestUp(t *testing.T) {
	s1 := "ChiHuo"

	fmt.Println(strings.ToLower(s1)) // 全小写

	s2 := "GOLANG_HHH"
	r := []rune(s2)
	if unicode.IsUpper(r[0]) {
		r[0] = unicode.ToLower(r[0])
	}

	fmt.Println(string(r)) // 首字母小写

}

// parsingUa 解析用户代理字符串，提取操作系统、WebKit版本和移动设备标识
func parsingUa(ua string) map[string]string {
	// 同时匹配多种结果 匹配到任意一个就是有结果
	matchRegex := map[string]*regexp.Regexp{
		"os":     regexp.MustCompile(`CPU iPhone OS ([0-9_])* like Mac OS X`), // 正则表达式用于匹配操作系统
		"webkit": regexp.MustCompile(`AppleWebKit.[\d\.]*`),                   // 正则表达式用于匹配WebKit版本
		"mobile": regexp.MustCompile(`Mobile\/([\dA-Z]*)`),                    // 正则表达式用于匹配移动设备标识
	}

	matchResult := make(map[string]string)

	for key, regex := range matchRegex {
		matches := regex.FindStringSubmatch(ua) // 查找匹配项
		// 匹配到才有第一个元素
		if len(matches) > 0 {
			matchResult[key] = matches[0]                     // 将匹配结果添加到结果映射中 mobile Mobile/16D57
			tmpKey := fmt.Sprintf("%s-%d", key, len(matches)) // 就算有多个也只匹配到第一个
			matchResult[tmpKey] = matches[1]                  // mobile-1 16D57
		}
	}

	return matchResult
}

func TestRegex(t *testing.T) {
	// 示例用户代理字符串
	// userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_1_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/16D57"

	userAgent := "Mozilla/5.0 Mobile/16D57"
	// userAgent := "Mobile/16D57 Mozilla/5.0 Mobile/16D58" // 就算有多个也只匹配到第一个
	// 解析用户代理
	result := parsingUa(userAgent)

	// 打印结果
	for key, value := range result {
		println(key, value)
	}
}

func TestRegex2(t *testing.T) {
	// str := "xxx-xxx-xxx-xxx-[xxx1]"
	str := "他趣-广点通-达人一口价-[非荷-[xcxx]尔蒙]-口播-[一手]-玉莹-xxxx"
	pattern := `-\[(.*)\]-?`
	// .* 非荷-[xcxx
	// .*? 非荷-[xcxx]尔蒙

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("正则表达式编译错误: %v\n", err)
		return
	}

	matches := re.FindStringSubmatch(str)
	if len(matches) > 1 {
		fmt.Printf("匹配到的内容: %s\n", matches[1])
	} else {
		fmt.Println("没有匹配到内容")
	}
}

// func MustCompile(str string) *Regexp {
// 	regexp, err := Compile(str)
// 	if err != nil {
// 		panic(`regexp: Compile(` + quote(str) + `): ` + err.Error())
// 	}
// 	return regexp
// }

func TestRegex3(t *testing.T) {
	text := "他趣-广点通-达人一口价-[非荷尔蒙]-口播-[一手]-玉莹-xxxx"

	// 正则表达式模式说明：
	// ()捕获组  反斜杠\  ?贪婪模式
	// [^-]只要中间不包含短横线的任意字符
	re := regexp.MustCompile(`([^-]+)-`)
	// [^-]+匹配1个或多个非短横线字符，比固定长度{3}更灵活

	// 查找所有匹配项（返回[][]string结构）
	matches := re.FindAllStringSubmatch(text, -1) // 使用FindAllStringSubmatch获取完整匹配上下文，避免边界问题
	fmt.Println(len(matches))
	for pos, match := range matches {
		// 确保捕获组有效性，避免空匹配
		if len(match) > 1 { // match[0]是完整匹配，match[1]是捕获组
			// fmt.Println(match[0]) // 达人一口价-
			fmt.Println(pos, match[1]) // 达人一口价
			// strings.Trim(match[1], "[")
		}
	}

}

func TestValuate(t *testing.T) {
	// 1. 创建表达式
	exprStr := "10 > 0 && (response_time <= 100 || retries < 5)"
	expression, err := govaluate.NewEvaluableExpression(exprStr)
	if err != nil {
		panic(err)
	}

	// 2. 准备参数
	params := map[string]interface{}{
		"response_time": 85.0,
		"retries":       3,
	}

	// 3. 计算表达式
	result, err := expression.Evaluate(params)
	if err != nil {
		panic(err)
	}

	// 4. 处理结果 (布尔值)
	fmt.Printf("结果: %v\n", result) // 输出: 结果: true

	// 函数调用
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return float64(length), nil
		},
	}
	expression2, _ := govaluate.NewEvaluableExpressionWithFunctions(
		"strlen('heo') > 3", // 比较符2边都必须是数值类型
		functions,
	)
	params2 := map[string]interface{}{} // 本身不需要参数的函数
	result2, err := expression2.Evaluate(params2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("结果: %v\n", result2)

	// exprStr := `
	// 	(employee_salary > 50000 &&
	// 	years_experience >= 5) ||
	// 	has_advanced_degree == true
	// `

	// params := map[string]interface{}{
	// 	"employee_salary":     65000.0,
	// 	"years_experience":    4,
	// 	"has_advanced_degree": true,
	// }

}

func TestValuate2(t *testing.T) {
	// 注册自定义函数
	functions := map[string]govaluate.ExpressionFunction{
		"max": func(args ...interface{}) (interface{}, error) {
			a := args[0].(float64)
			b := args[1].(float64)
			return math.Max(a, b), nil
		},
	}

	// 创建带函数的表达式
	expr, _ := govaluate.NewEvaluableExpressionWithFunctions(
		"max(temperature, 30) > 25",
		functions,
	)

	params := map[string]interface{}{"temperature": 28.5}
	result, _ := expr.Evaluate(params)
	fmt.Printf("结果: %v\n", result)
}

// 数学运算 "(price * quantity) + tax"
// 逻辑运算 "score >= 80 && attendance > 0.75"
// 三元表达式 "status == 200 ? 'OK' : 'ERROR'"
// 函数调用
