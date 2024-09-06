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
