package context_cancel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func isCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done(): //接收取消通知
		return true
	default:
		return false
	}
}

func TestCancel(t *testing.T) {

	//当context被取消时 对应的子context也会取消
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 5; i++ {
		go func(i int, ctx context.Context) {
			for {
				if isCancelled(ctx) {
					break
				}
				time.Sleep(time.Millisecond * 5)
			}
			fmt.Println(i, "Cancelled")
		}(i, ctx)
	}
	//取消根节点
	cancel()
	time.Sleep(time.Second * 1) // 等待子的执行完毕
	fmt.Println("end")
}

func TestCancel2(t *testing.T) {

	//当context被取消时 对应的子context也会取消
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int, ctx context.Context, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				if isCancelled(ctx) {
					break // 收到通知跳出循环
				}
				time.Sleep(time.Millisecond * 5)
			}
			fmt.Println(i, "Cancelled")

		}(i, ctx, &wg)
	}
	//取消根节点
	cancel()
	wg.Wait() // 等待子的执行完毕
	fmt.Println("end")
}
