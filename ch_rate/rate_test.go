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
	limit := rate.Limit(5)                   // QPS：50 基础速率
	burst := 25                              // 桶容量25
	limiter := rate.NewLimiter(limit, burst) // 1. 初始化一个令牌生成速率为limit，容量为burst的令牌桶

	size := 50 // 数据量500
	data := generateData(size)

	var wg sync.WaitGroup // 工作组锁
	startTime := time.Now()
	for i, item := range data {
		wg.Add(1)
		go func(idx int, obj int) {
			defer wg.Done()
			// 2. Wait拿到令牌
			if err := limiter.Wait(context.Background()); err != nil {
				t.Logf("[%d] [EXCEPTION] wait err: %v", idx, err)
			}
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

func TestTimeRate(t *testing.T) {

	limiter := rate.NewLimiter(10, 100) // 1. 初始化一个令牌生成速率为limit，容量为burst的令牌桶

	timer := time.NewTimer(time.Second * 10) // 定时器
	quit := make(chan struct{})              // 通道
	defer timer.Stop()
	go func() {
		<-timer.C   // 定时器到期
		close(quit) // 通知子协程都关闭
	}()

	var allowed, denied int32
	var wait sync.WaitGroup
	cpuNum := runtime.NumCPU()
	fmt.Println(cpuNum)
	for i := 0; i < cpuNum; i++ {
		wait.Add(1)
		go func() {
			for {
				select {
				case <-quit:
					wait.Done()
					return
				default:
					ctx, cancel := context.WithTimeout(context.TODO(), time.Second) // 1s超时
					// 2. Wait拿到令牌
					err := limiter.Wait(ctx)
					if err == nil {
						atomic.AddInt32(&allowed, 1)
					} else {
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
