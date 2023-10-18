package concurrency

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func runTask(id int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("The result is from %d", id) //只是返回字符串
}

func FirstResponse() string {
	numOfRunner := 10
	//ch := make(chan string) //如果不用缓存的 会导致其他chan等待接收的协程阻塞，如果多了会导致系统资源耗尽
	ch := make(chan string, numOfRunner) //防止协程泄露

	for i := 0; i < numOfRunner; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret //数据写入channel
		}(i)
	}
	//有缓存 的获取chan数据
	return <-ch
}

func TestFirstResponse(t *testing.T) {
	t.Log("Before:", runtime.NumGoroutine())
	t.Log(FirstResponse())
	time.Sleep(time.Second * 1)
	t.Log("After:", runtime.NumGoroutine())

}
