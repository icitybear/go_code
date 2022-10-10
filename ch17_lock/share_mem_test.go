package share_mem

import (
	"sync"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {

	counter := 0
	for i := 0; i < 5000; i++ {
		//每起一个协程 共享的都是同个栈区的局部变量 因为无序的 共享内存了 所以counter 到不了 5000
		go func() {
			counter++
		}() //无序 导致计数不准
	}
	time.Sleep(1 * time.Second)
	t.Logf("counter = %d", counter)

}

func TestCounterThreadSafe(t *testing.T) {
	var mut sync.Mutex //锁（Mutex完全互斥锁）RWLock 读写锁
	counter := 0
	for i := 0; i < 5000; i++ {
		go func() {
			defer func() {
				mut.Unlock() //解锁
			}()
			mut.Lock() //后续的协程解锁了 才继续计数 这样才有序
			counter++  //非值传递 共享内存
		}()
	}
	time.Sleep(1 * time.Second) //加这个是因为 如果外面执行完的速度 超过了所有协程执行完 这样输出就不准确
	t.Logf("counter = %d", counter)

}

func TestCounterWaitGroup(t *testing.T) {
	var mut sync.Mutex
	var wg sync.WaitGroup //等待锁
	counter := 0
	for i := 0; i < 5000; i++ {

		wg.Add(1) //等待加1   wg.Wait()有了这个 外面就不用time.Sleep了

		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock() // 不加锁 并发可能会有多次的 10++ 保持精准却耗性能
			counter++  //counter是共享内存 同一个栈区地址

			wg.Done() //等待完成1个
		}()
	}
	wg.Wait() //等待中 就是等锁都执行完了 才往后继续执行 代替time.Sleep了
	t.Logf("counter = %d", counter)

}
