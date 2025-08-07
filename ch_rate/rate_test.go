package pinyin_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"
) // 需要import的rate库，其它import暂时忽略

// 生成0->X的数据集
func generateData(num int) []int {
	var data []int
	for i := 0; i < num; i++ {
		data = append(data, i)
	}
	return data
}

// 处理数据，数字*10
func process(obj interface{}) (interface{}, error) {
	integer, ok := obj.(int)
	if !ok {
		return nil, errors.New("invalid integer")
	}
	time.Sleep(1)
	nextInteger := integer * 10
	if integer%99 == 0 {
		return nextInteger, errors.New("not a happy number")
	}
	return nextInteger, nil
}

func TestRate(t *testing.T) {
	// num := rate.Every(time.Millisecond * 5)
	num := 5
	limit := rate.Limit(num)                 // QPS：50 基础速率
	burst := 300                             // 桶容量25
	limiter := rate.NewLimiter(limit, burst) // 1. 初始化一个令牌生成速率为limit，容量为burst的令牌桶

	size := 500 // 数据量500
	data := generateData(size)

	var wg sync.WaitGroup // 工作组锁
	startTime := time.Now()
	for i, item := range data {
		wg.Add(1)
		go func(idx int, obj int) {
			defer wg.Done()
			t.Logf("idx:%d, start: %v", idx, time.Now().Format("2006-01-02 15:04:05.000"))
			// 2. Wait拿到令牌
			if err := limiter.Wait(context.Background()); err != nil {
				t.Logf("[%d] [EXCEPTION] wait err: %v", idx, err)
			}
			t.Logf("idx:%d, End: %v", idx, time.Now().Format("2006-01-02 15:04:05.000"))
			// 执行业务逻辑
			processed, err := process(obj)
			if err != nil {
				// 也要模拟处理时的报错
				t.Logf("[%d] [ERROR] processed: %v, err: %v", idx, processed, err)
			} else {
				t.Logf("[%d] [OK] processed: %v", idx, processed)
			}
		}(i, item)
	}
	wg.Wait()
	endTime := time.Now()
	t.Logf("start: %v, end: %v, seconds: %v", startTime, endTime, endTime.Sub(startTime).Seconds())
}

func TestTimeRate0(t *testing.T) {

	num := rate.Every(time.Millisecond * 100) // 0.1=>每秒10个 100写成500每秒50个
	// num := rate.Limit(10)                 // 10（每秒10个） 20 => 0.2  等价
	limiter := rate.NewLimiter(num, 1000) // 1. 初始化一个令牌生成速率为limit，容量为burst的令牌桶
	// limiter := rate.NewLimiter(0.05, 2)
	fmt.Println(limiter.Burst()) // 1000
	fmt.Println(limiter.Limit()) // 1/0.1 = 10

}

func TestTimeRate(t *testing.T) {

	num := rate.Every(time.Millisecond * 100) // 0.1=>每秒10个 100写成500每秒50个
	// num := rate.Limit(10)                // 10=>0.1（每秒10个） 20 => 0.2
	limiter := rate.NewLimiter(num, 5) // 1. 初始化一个令牌生成速率为limit，容量为burst的令牌桶

	// tag:容量为burst容量一开始就满了，然后下个时间频率继续产令牌
	timer := time.NewTimer(time.Second * 3) // 定时器 5秒的定时器 5s后通道收到消息
	quit := make(chan struct{})             // 通道
	defer timer.Stop()
	// tag:单独起一个协程 监听定时器收到消息主动close通道
	go func() {
		<-timer.C   // 定时器到期 定时器收到消息
		close(quit) // 通知子协程都关闭
	}()

	var allowed, denied int32
	var wait sync.WaitGroup
	cpuNum := runtime.NumCPU()
	fmt.Println(cpuNum) // 8核

	for i := 0; i < cpuNum; i++ {
		wait.Add(1)
		go func() {
			// for + select 等价 range 持续监听
			for {
				select {
				// 接收外部关闭消息
				case <-quit:
					wait.Done()
					return
				default:
					ctx, cancel := context.WithTimeout(context.TODO(), time.Second) // 1s超时
					// 2. Wait拿到令牌
					err := limiter.Wait(ctx)
					if err == nil {
						t.Logf("获取到令牌: allowed %v, Time: %v", allowed, time.Now().Format("2006-01-02 15:04:05.000")) //打印每个抢到的时间
						atomic.AddInt32(&allowed, 1)
					} else {
						// ctx设置1秒超时 1s都有1万次 rate: Wait(n=1) would exceed context deadline
						// 如果令牌产生太慢就会有N个这日志
						t.Logf("未获取到令牌:denied %v, Time: %v, err: %s", denied, time.Now().Format("2006-01-02 15:04:05.000"), err.Error())
						fmt.Println(err)
						atomic.AddInt32(&denied, 1)
					}
					cancel()
				}
			}
		}()
	}

	wait.Wait()
	fmt.Printf("allowed: %d, denied: %d, qps: %d\n", allowed, denied, (allowed+denied)/10)
}
