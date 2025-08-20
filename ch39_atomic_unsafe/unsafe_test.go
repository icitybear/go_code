package unsafe_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
)

type Customer struct {
	Name string
	Age  int
}

func TestUnsafe(t *testing.T) {
	i := 10
	f := *(*float64)(unsafe.Pointer(&i))
	t.Log(unsafe.Pointer(&i))
	t.Log(f)
}

// The cases is suitable for unsafe
type MyInt int

// 合理的类型转换
func TestConvert(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := *(*[]MyInt)(unsafe.Pointer(&a))
	t.Log(b)
}

// 原子类型操作
func TestAtomic(t *testing.T) {
	var shareBufPtr unsafe.Pointer
	writeDataFn := func() {
		data := []int{}
		for i := 0; i < 100; i++ {
			data = append(data, i)
		}
		atomic.StorePointer(&shareBufPtr, unsafe.Pointer(&data))
	}
	readDataFn := func() {
		data := atomic.LoadPointer(&shareBufPtr)
		fmt.Println(data, *(*[]int)(data))
	}
	var wg sync.WaitGroup
	writeDataFn()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				writeDataFn()
				time.Sleep(time.Microsecond * 100)
			}
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				readDataFn()
				time.Sleep(time.Microsecond * 100)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

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

func TestAtomicCsp(t *testing.T) {
	c = 0
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go worker2(&wg)
	}
	wg.Wait()
	fmt.Printf("c = %d", c)
}

// 使用并发编程提高响应的案例
func processDeal(val int) int {
	time.Sleep(1 * time.Second) // 耗时久一点，协程执行时的并发才能体现
	return val * 2
}

// 原子 错误用法 正确的用法
func TestCspErr(t *testing.T) {
	start := time.Now()

	arrSlice := []int{1, 2, 3, 4, 5}
	var result []int

	wg := sync.WaitGroup{}
	// var atomicVal atomic.Value // =》 any  整个对象存储
	// atomicVal.Store(make([]int, 0)) // 一开始就指定类型
	// if v := atomicVal.Load(); v != nil {
	// 	atomicVal.Store(result)
	// }

	var atomicVal int32
	for pos, v := range arrSlice {
		wg.Add(1)
		_ = pos
		go func(val int) {
			defer wg.Done()
			newVal := processDeal(val)
			// 使用cas保持原子性  错误用法
			old := atomic.LoadInt32(&atomicVal) // 如果多个goroutine读取到相同的old值，就会导致部分CAS失败
			// CAS 仅保护了计数器 atomicVal 的原子递增
			if atomic.CompareAndSwapInt32(&atomicVal, old, old+1) {
				// 即使CAS成功，多个goroutine同时append同一个切片也是不安全的。
				result = append(result, newVal)
			} else {
				fmt.Println("cas原子操作失败", old)
			}
		}(v)
	}
	wg.Wait()
	spew.Printf("result:%+v \n", result)
	spew.Printf("atomicVal:%+v \n", atomic.LoadInt32(&atomicVal))
	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}

// result:[6 4] // 错误情况
// atomicVal:5	// 值可能非5 4
// 程序执行耗时(s): 1.001535792

// 比如3个协程跑
//	Time | Goroutine A         | Goroutine B         | Goroutine C
//
// -----|---------------------|---------------------|--------------------
//
//	t1  | Load old=0          |                     |
//	t2  |                     | Load old=0          |
//	t3  | CAS(0→1) 成功       |                     |
//	    | append(result)      |                     |
//	t4  |                     | CAS(0→1) 失败!      |
//	t5  |                     |                     | Load old=1
//	t6  |                     |                     | CAS(1→2) 成功
//	    |                     |                     | append(result)
//
