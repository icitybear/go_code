package cancel

import (
	"context"
	"fmt"
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
	time.Sleep(time.Second * 1)
}
