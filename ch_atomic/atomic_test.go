package atomic_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

var c int64

// 没用锁 可能资源竞争问题 不准确
func worker1(wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		c++
	}
	wg.Done()
}

// worker1和worker2区别 无锁和有原子操作atomic
// 锁针对的是代码块-过程  atomic针对的是变量-硬件 cpu指令
// 主要是 Add、CompareAndSwap、Load、Store、Swap。 https://blog.csdn.net/h_l_f/article/details/118739317
func worker2(wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		// 传入变量指针 加多少
		atomic.AddInt64(&c, 1)
	}
	wg.Done()
}

func TestAtomic(t *testing.T) {
	c = 0
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go worker2(&wg)
	}
	wg.Wait()
	fmt.Printf("c = %d", c)
}
