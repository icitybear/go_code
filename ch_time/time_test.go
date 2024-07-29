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

func TestBet(t *testing.T) {
	startDateTime, _ := time.Parse("2006-01-02", "2024-06-04")
	endDateTime, _ := time.Parse("2006-01-02", "2024-06-05")
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

func TestFormat(t *testing.T) {

	a := 1632 % 1000
	fmt.Println(a)

	statDateTime := time.Now().AddDate(0, 0, -1)
	str := statDateTime.Format("2006-01-02 00:00:00")
	fmt.Println(str)

	str1 := "2024-06-07" // 默认是0点
	parsedTime1, _ := time.ParseInLocation("2006-01-02", str1, time.Local)
	parsedTime2 := time.Now().AddDate(0, 0, -14)
	// After 这里直接对比的是日期 Ymd没包括时分秒
	if parsedTime1.After(parsedTime2) {
		fmt.Println(str1)
		return
	}
	fmt.Println("after")
}
