package goroutine_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func dosomething(o *sync.Once) {
	fmt.Println("Start:")
	o.Do(func() {
		fmt.Println("Do Something...")
	})
	fmt.Println("Finished.")
}

func TestOnce(t *testing.T) {
	o := &sync.Once{}
	go dosomething(o)
	go dosomething(o)
	time.Sleep(time.Second * 1)
}

// 输出
// Start:
// Do Something.. 只初始化执行一次
// Finished.
// Start:
// Finished.
