package time_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type CronParam struct {
	Key   string
	Value int
	Every bool
}

func GenerateCron(params []CronParam) string {
	second := "0"
	minute := "0"
	hour := "0"
	day := "*"
	month := "*"
	week := "*"

	hourFirst := false
	for _, p := range params {
		switch p.Key {
		case "month":
			if p.Every {
				month = fmt.Sprintf("*/%d", p.Value)
			} else {
				month = appendOrSet(month, p.Value)
			}
		case "day":
			if p.Every {
				day = fmt.Sprintf("*/%d", p.Value)
			} else {
				day = appendOrSet(day, p.Value)
			}
		case "week":
			// 周没有间隔的写法 0-6
			if !p.Every {
				week = appendOrSet(week, p.Value)
			}
		case "hour":
			if p.Every {
				hour = fmt.Sprintf("*/%d", p.Value)
			} else {
				if !hourFirst {
					hourFirst = true
					hour = fmt.Sprintf("%d", p.Value)
				} else {
					hour = appendOrSet(hour, p.Value)
				}
			}
		}
	}

	// 处理 Day 和 Week 冲突 都有值的情况
	if day != "*" && week != "*" {
		// if day != "?" && week != "?" {
		// 	day = "?"
		// }
		week = "*"
	}

	return fmt.Sprintf("%s %s %s %s %s %s", second, minute, hour, day, month, week)
}

func appendOrSet(field string, value int) string {
	if field == "*" || field == "?" {
		return strconv.Itoa(value)
	}
	return fmt.Sprintf("%s,%d", field, value)
}

// 示例调用
func TestCron(t *testing.T) {
	// 示例1：1月的3、4、13号
	params1 := []CronParam{
		{Key: "month", Value: 1, Every: false},
		{Key: "day", Value: 3, Every: false},
		{Key: "day", Value: 4, Every: false},
		{Key: "day", Value: 13, Every: false},
	}
	fmt.Println(GenerateCron(params1)) // 输出: 0 0 0 3,4,13 1 * *

	// 同级别设定冲突的时候就是以最后一个为准
	// 每1月的3号 5号
	params2 := []CronParam{
		{Key: "month", Value: 1, Every: true},
		{Key: "day", Value: 3, Every: false},
		{Key: "day", Value: 5, Every: false},
	}
	fmt.Println(GenerateCron(params2)) // 0 0 0 3,5 */1 *

	// 示例3：每周三、周四的8点
	params3 := []CronParam{
		{Key: "week", Value: 3, Every: false},
		{Key: "week", Value: 4, Every: false},
		{Key: "hour", Value: 8, Every: false},
	}
	fmt.Println(GenerateCron(params3)) // 输出: 0 0 8 * * 3,4

	params4 := []CronParam{
		{Key: "day", Value: 6, Every: false},
		{Key: "day", Value: 28, Every: false},
		{Key: "month", Value: 1, Every: true},
		{Key: "week", Value: 0, Every: false},
		{Key: "week", Value: 2, Every: false},
		{Key: "hour", Value: 3, Every: true},
	}
	// 虽然日和周冲突了 默认日为准 0 0 */3 6,28 */1 *
	fmt.Println(GenerateCron(params4)) // 0 0 */3 ? */1 0,2

	params5 := []CronParam{
		// {Key: "week", Value: 1, Every: true}, // 如果周重复了只需要传一个
		{Key: "week", Value: 5, Every: false},
		{Key: "week", Value: 0, Every: false},
		{Key: "week", Value: 1, Every: false},
		{Key: "hour", Value: 2, Every: false},
	}
	fmt.Println(GenerateCron(params5)) // 0 0 0,2 * * 5,0,1
}

func ParseCron(cron string) []CronParam {
	parts := strings.Split(cron, " ")
	if len(parts) != 6 {
		return nil
	}

	hour := parts[2]
	day := parts[3]
	month := parts[4]
	week := parts[5]

	var params []CronParam

	// 解析小时（hour）
	if hourParams := parseField(hour, "hour"); hourParams != nil {
		for _, p := range hourParams {
			if p.Value != 0 { // 忽略 0 点（如无特殊需求）
				params = append(params, p)
			}
		}
	}

	// 解析日期（day）
	params = append(params, parseField(day, "day")...)

	// 解析月份（month）
	params = append(params, parseField(month, "month")...)

	// 解析星期（week）
	params = append(params, parseField(week, "week")...)

	return params
}

// 解析单个字段，返回对应的 CronParam 列表
func parseField(field, key string) []CronParam {
	var params []CronParam

	if field == "*" || field == "?" {
		return params
	}

	// 处理间隔值 */N
	if strings.HasPrefix(field, "*/") {
		if value, err := strconv.Atoi(field[2:]); err == nil {
			params = append(params, CronParam{Key: key, Value: value, Every: true})
		}
		return params
	}

	// 处理逗号分隔的固定值列表
	values := strings.Split(field, ",")
	for _, v := range values {
		if value, err := strconv.Atoi(v); err == nil {
			params = append(params, CronParam{Key: key, Value: value, Every: false})
		}
	}

	return params
}

// 示例调用
func TestCron2(t *testing.T) {
	// 示例1：0 0 0 3,4,13 1 * *
	cron1 := "0 0 0 3,4,13 1 * *"
	fmt.Printf("%#v\n", ParseCron(cron1))
	// 输出: []main.CronParam{
	//   {Key:"day", Value:3, Every:false},
	//   {Key:"day", Value:4, Every:false},
	//   {Key:"day", Value:13, Every:false},
	//   {Key:"month", Value:1, Every:false}
	// }

	// 示例2：0 0 0 */3 1 * *
	cron2 := "0 0 0 */3 1 * *"
	fmt.Printf("%#v\n", ParseCron(cron2))
	// 输出: []main.CronParam{
	//   {Key:"day", Value:3, Every:true},
	//   {Key:"month", Value:1, Every:false}
	// }

	// 示例3：0 0 8 * * 3,4
	cron3 := "0 0 8 * * 3,4"
	fmt.Printf("%#v\n", ParseCron(cron3))
	// 输出: []main.CronParam{
	//   {Key:"hour", Value:8, Every:false},
	//   {Key:"week", Value:3, Every:false},
	//   {Key:"week", Value:4, Every:false}
	// }
}
