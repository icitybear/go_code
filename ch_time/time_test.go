package time_test

import (
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// 时间字符串格式校验
	// dtTime, err := time.Parse("2006-01-02", "2024/02/03")
	// fmt.Println(dtTime, err)

	t1, err := time.Parse("2006-01-02", "2024-05-24")
	fmt.Println(t1) // 2024-05-24 00:00:00 +0000 UTC
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(t1.Unix())         // 1716508800  2024-05-24 08:00:00
	fmt.Println(t1.Local().Unix()) // 1716508800 // UTC时区

	t2 := time.Now().Unix()
	fmt.Println(t2) // 1716625124

	t3, _ := time.ParseInLocation("2006-01-02", "2024-05-24", time.Local)
	fmt.Println(t3)        // 2024-05-24 00:00:00 +0800 CST
	fmt.Println(t3.Unix()) // 1716480000
}

func TestParseHour(t *testing.T) {
	loc := time.Local
	str := "2025-06-22 1:00" // hour=1
	tsTime, err := time.ParseInLocation("2006-01-02 15:04", str, time.Local)
	if err != nil {
		fmt.Println(err)
	}
	statDate := tsTime.Format("20060102")
	statHour := tsTime.Hour()
	fmt.Println(statDate, statHour)

	// 自然日
	res := time.Now().Sub(tsTime)
	fmt.Println(res)

	start := time.Date(tsTime.Year(), tsTime.Month(), tsTime.Day(), 0, 0, 0, 0, loc)
	// 规范化到结束日的 00:00:00
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// 计算时间差并转换为天数
	fmt.Println(end.Sub(start).Hours())
	days := end.Sub(start).Hours() / 24
	fmt.Println(days) //  如果要算自然日就要+1

}

func TestParse2(t *testing.T) {
	// 最后相同值会覆盖覆盖
	// str := "2025-05-25 00:00 - 00:59"
	str := "2025-05-25 00:00 - 01:59" // hour=1
	tsTime, err := time.ParseInLocation("2006-01-02 15:04 - 15:04", str, time.Local)

	if err != nil {
		fmt.Println(err)
	}
	statDate := tsTime.Format("20060102")
	statHour := tsTime.Hour()
	fmt.Println(statDate, statHour)
}
func TestBet(t *testing.T) {
	startDateTime, _ := time.Parse("2006-01-02", "2024-12-30")
	endDateTime, _ := time.Parse("2006-01-02", "2025-01-02")
	i := 0
	for {
		statDateTime := startDateTime.AddDate(0, 0, i)
		// 不包含结束时间 只打印了2024-06-04
		if !statDateTime.Before(endDateTime) {
			// fmt.Println(statDateTime)
			return
		}
		i++

		statDate := statDateTime.Format("2006-01-02")
		fmt.Println(statDate)

	}
}

func TestFormatDb(t *testing.T) {

	parsedTime := time.Now().AddDate(0, 0, -7)
	fmt.Println(parsedTime.Format("2006-01-02"))
	parsedTime1, _ := time.ParseInLocation("2006-01-02", "2025-06-07", time.Local)
	parsedTime2, _ := time.ParseInLocation("2006-01-02", "2025-06-08", time.Local)
	// tag: before after不包含当日
	if parsedTime1.Before(parsedTime2) {
		fmt.Println("before")
	}
	// After 这里直接对比的是日期 Ymd没包括时分秒
	if parsedTime1.After(parsedTime2) {
		fmt.Println("after")
	}
	fmt.Println("end")
}

func TestFormatT(t *testing.T) {

	loc, _ := time.LoadLocation("Local") // 不能写成local小写 会报错

	t1, err := time.ParseInLocation("2006-01-02", "2024-11-07", loc)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t1.Unix())

	// 将时间戳转换为 time.Time 类型
	timestamp := int64(1731398862)
	t2 := time.Unix(timestamp, 0)

	// 将时间转换为指定格式的日期字符串
	dateString := t2.Format("2006-01-02 15:03:04")

	fmt.Println(dateString, t2.Unix())

	endTime := time.Now().AddDate(0, 0, -1)
	fmt.Println(endTime.Unix())
}

func TestXXX(t *testing.T) {
	eventTime := 1745764920
	//regTime := 1692860732
	reflowTime := 1745755798
	res := IsSameDay(int64(eventTime), int64(reflowTime))
	fmt.Println(res)

}

func IsSameDay(aTime int64, bTime int64) bool {
	deviceY, deviceM, deviceD := time.Unix(aTime, 0).Date()
	accountY, accountM, accountD := time.Unix(bTime, 0).Date()
	if deviceY == accountY && deviceM == accountM && deviceD == accountD {
		return true
	}

	return false
}

func TimestampToDateString(timestamp int64) string {
	// loc, _ := time.LoadLocation("Local")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// 将时间戳转换为 time.Time 类型 .UTC() 容器里linux时区默认UTC
	t := time.Unix(timestamp, 0).In(loc)

	// 将时间转换为指定格式的日期字符串
	dateString := t.Format("2006-01-02 15:04:05")
	return dateString
}
func TestFormatT2(t *testing.T) {
	ts := 1743673260
	str := TimestampToDateString(int64(ts))
	println(str)
}
func TestToday(t *testing.T) {
	startDate := "2025-02-11"
	// 解析日期字符串
	parsedDate, _ := time.Parse("2006-01-02", startDate)
	// 获取当前日期（去掉时间部分）
	today := time.Now().Truncate(24 * time.Hour) // 因为这种方式无法准确地将时间截取到当天的日期。
	// 正确的方法是通过比较日期的年、月、日部分来判断是否为同一天。
	if parsedDate == today {
		fmt.Println(today)
	} else {
		fmt.Println("非当日")
	}
	now := time.Now()
	if parsedDate.Year() == now.Year() &&
		parsedDate.Month() == now.Month() &&
		parsedDate.Day() == now.Day() {
		fmt.Println(now)
	} else {
		fmt.Println("非当日2")
	}
}
func TestFormatTrans(t *testing.T) {

	statDateTime := time.Now()
	str := statDateTime.Format("2006/01/02 00:00:00")
	fmt.Println(str)

	// 先转瞬逝时间再使用格式化
	str1 := "2024-06-07"                                                   // 默认是0点
	parsedTime1, _ := time.ParseInLocation("2006-01-02", str1, time.Local) // 2个参数都格式都要一一对上 0001-01-01 00:00:00 +0000 UTC
	transStr := parsedTime1.Format("20060102")
	fmt.Println(transStr)
}

func TestTimer(t *testing.T) {
	ch := make(chan int)
	// 起协程
	go func() {
		// 1. for + select 持续监听
		for {
			// 2. select可以完成监控多个channel的状态, 不同channel 收到消息执行顺序与case无关 如果都没收到 就默认走default
			select {
			case num := <-ch: // 外层是1秒就发一次 无缓冲区
				fmt.Println("get num is ", num)
			case <-time.After(2 * time.Second): // 每次都是新的2秒定时。
				// 3. 每次调用 time.After 都会创建一个新的计时器, 底层的计时器在计时器被触发之前不会被垃圾收集器回收。存在内存泄漏的可能
				fmt.Println("time's up!!!")
				// 如果都没收到 就默认走default 如果没有 default 子句,select 将阻塞,直到某个通道可以运行;
			}
		}
	}()

	for i := 0; i < 5; i++ {
		ch <- i
		time.Sleep(1 * time.Second)
	}
}

func TestBefore(t *testing.T) {
	startDateTime, _ := time.Parse("2006-01-02", "2025-06-01")
	endDateTime, _ := time.Parse("2006-01-02", "2025-06-09")
	// 开始日期
	statDt := startDateTime
	i := 5 // 夸天执行
	for {

		endDt := statDt.AddDate(0, 0, i)
		if statDt.After(endDateTime) {
			fmt.Print(statDt, "limit")
			break
		}

		if endDt.After(endDateTime) {
			endDt = endDateTime
		}

		statDate := statDt.Format("2006-01-02")
		endDate := endDt.Format("2006-01-02")
		fmt.Println(statDate, endDate)

		statDt = endDt.AddDate(0, 0, 1)
	}
}

func TestFunc(t *testing.T) {
	ttl := RemainingTimeNDay(time.Now().Unix(), 1) // 1表示当日
	fmt.Println(ttl)

	startTime := time.Now()
	time.Sleep(3 * time.Second)
	fmt.Println(time.Since(startTime).Seconds())
}

// RemainingTimeNDay 计算n天后的剩余时间戳
func RemainingTimeNDay(t1 int64, n int) int64 {
	// 转为时间
	t := time.Unix(t1, 0)

	// 获取今天的结束时间
	endTime := time.Now().AddDate(0, 0, n-1)
	endOfDay := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, time.Now().Location())

	// 计算剩余时间
	remainingTime := endOfDay.Sub(t)

	// 将剩余时间转换为秒，并返回
	return int64(remainingTime.Seconds())
}

// time.Sleep 改成 time.After 可以降低 cpu 使用率
// time.Sleep 不会让度 cpu 使用权（在这个代码里，cpu 被白白浪费 100ms 的时间）
// time.After 会让度出来

// GetISOWeekNumber 根据周一和周日的日期计算ISO周数和年份
// 参数: mondayStr, sundayStr - 格式为"2006-01-02"的日期字符串
// 返回值: year, week - ISO年份和周数；error - 错误信息
func GetISOWeekNumber(mondayStr, sundayStr string) (year, week int, err error) {
	// 解析日期字符串（使用UTC时区）
	layout := "2006-01-02"
	monday, err := time.Parse(layout, mondayStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid Monday date: %v", err)
	}
	sunday, err := time.Parse(layout, sundayStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid Sunday date: %v", err)
	}

	// 验证周一和周日的关系
	if monday.Weekday() != time.Monday {
		return 0, 0, fmt.Errorf("the first date is not a Monday")
	}
	if sunday.Weekday() != time.Sunday {
		return 0, 0, fmt.Errorf("the second date is not a Sunday")
	}
	if sunday.Sub(monday) != 6*24*time.Hour {
		return 0, 0, fmt.Errorf("the interval between dates is not exactly 6 days")
	}

	// 使用周四计算ISO周（确保跨年周计算正确）
	thursday := monday.Add(3 * 24 * time.Hour)
	year, week = thursday.ISOWeek()
	return year, week, nil
}

// 示例使用
func TestFuncWeek(t *testing.T) {
	// 示例1: 2025年第一周（跨年）
	monday := "2024-12-30"
	sunday := "2025-01-05"
	year, week, err := GetISOWeekNumber(monday, sunday)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Date: %s to %s -> ISO Year: %d, Week: %d\n", monday, sunday, year, week)
	}

	monday = "2025-12-01"
	sunday = "2025-12-07"
	year, week, err = GetISOWeekNumber(monday, sunday)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Date: %s to %s -> ISO Year: %d, Week: %d\n", monday, sunday, year, week)
	}
}
