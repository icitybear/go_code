package context_cancel

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 可以看源码解析的那文章和视频
func TestContext(t *testing.T) {
	worker1()
}

// 自动超时
func worker1() {
	deep := 10
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	go handle(ctx, 500*time.Millisecond, deep)
	select {
	case <-ctx.Done():
		fmt.Println("worker1", ctx.Err())
	}
}

// 传值
func worker2() {
	deep := 10
	// 手动超时WithCancel
	ctx, cancel := context.WithCancel(context.Background())
	go handle(ctx, 500*time.Millisecond, deep)
	time.Sleep(1 * time.Second)
	cancel()
}

func worker3() {
	deep := 10
	// 装饰者 WithValue传值 key => value
	ctx := context.WithValue(context.Background(), "token", "citybear")
	// 自动超时WithTimeout 5秒 发给ctx done了
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// 发布结束
	defer cancel()

	// 定了0.5秒后 就使用select
	go handle(ctx, 500*time.Millisecond, deep)
	select {
	case <-ctx.Done():
		fmt.Println("worker3 ctx done", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration, deep int) {
	if deep > 0 {
		time.Sleep(200 * time.Millisecond)
		// 模拟自己调自己  上下文
		go handle(ctx, duration, deep-1)
	}
	fmt.Printf("init deep is %d \n", deep)
	// 获取传的值
	if ctx.Value("token") != nil {
		fmt.Printf("token is %s\n", ctx.Value("token"))
	}

	select {
	case <-ctx.Done():
		fmt.Println("handle ctx.done", ctx.Err())
		// 超时了 就执行
	case <-time.After(duration):
		fmt.Printf("process request with %v, %d\n", duration, deep)
	}
}
