package goroutine_test

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// 开启协程时的传参
func TestGroutine(t *testing.T) {

	for i := 0; i < 20; i++ {

		// i 值传递（拷贝值） 协程里的变量地址是不一样的 无竞争关系 就不会有共享内存了
		go func(i int) {
			//time.Sleep(time.Second * 1) //外面程序执行时间只要大于 协程，就会有输出了
			fmt.Println(i) // 协程执行打印输出也是无序的
		}(i)

		//竞争关系 运行时间快 导致输出的全是20 共享内存了
		// go func() {
		// 	fmt.Println(i) // 打印的时候i已经循环到20了
		// }()

		// 竞争关系 尤其是for循环结构体指针
	}

	time.Sleep(time.Millisecond * 100) //加这个是因为 外面执行完的速度 超过了协程 这样就没输出了
}

// time.Sleep 太low 如何保证子goroutine执行完 才结束主goroutine
var counter int = 0 // 自定义计数器 共享内存 通过锁保护

func addV2(a, b int, lock *sync.Mutex) {
	lock.Lock()
	c := a + b
	// 由于 counter 变量会被所有协程共享，为了避免 counter 值被污染 入了锁机制，即 sync.Mutex
	counter++
	fmt.Printf("%d: %d + %d = %d\n", counter, a, b, c)
	lock.Unlock()
}

// 方案1 共享内存变量+锁 = 计数器
func TestGroutine2(t *testing.T) {
	start := time.Now()
	lock := &sync.Mutex{}

	// 每执行一次子协程，该计数器的值加 1，当所有子协程执行完毕后，计数器的值应该是 10，
	for i := 0; i < 10; i++ {
		go addV2(1, i, lock) // 单机锁保证让goroutine并发争抢资源变成串行的
	}

	// 我们在主协程中通过一个死循环来判断 counter 的值，只有当它大于等于 10 时，才退出循环，进而退出整个程序
	for {
		lock.Lock()
		c := counter // 访问 其他goroutine的共享内存时 也要通过锁
		lock.Unlock()
		runtime.Gosched() // 让出 CPU 时间片 去执行子协程的  终止当前协程 runtime.Goexit() 阻塞当前
		if c >= 10 {
			break
		}
	}
	fmt.Printf("counter: %d\n", counter)
	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}

// 方案2 系统实现的工作组（封装好的计数器）
// wg.Done()为要执行的doneFunc() 外部计数加，子goroutine减
func addV3(a, b int, doneFunc func()) {
	defer func() {
		doneFunc() // 子协程执行完毕后将计数器-1
	}()
	c := a + b
	fmt.Printf("%d + %d = %d\n", a, b, c)
}

// 有效地管理并发任务的执行顺序和同步，确保所有的goroutine都执行完成后再进行下一步操作。
func TestGroutine3(t *testing.T) {
	start := time.Now()
	// 改为等待锁
	wg := sync.WaitGroup{}
	//wg.Add(10) // 初始化计数器数目为10 这样顺序还是乱的
	for i := 0; i < 5; i++ {
		wg.Add(1) // 运行时才能知道计数器的数目, 循环体内动态增加计数器，每次+1 顺序正确
		go addV3(1, i, wg.Done)
	}

	wg.Wait() // 等待子协程全部执行完毕退出  直到计数器归零
	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}

// 将线程加共享内存的方式称为「共享内存系统」,「消息传递系统」解决共享内存系统存在的问题(加锁来避免死锁或资源竞争)
// Communicating Sequential Processes 中提出的，在 CSP 系统中，所有的并发操作都是通过独立线程以异步运行的方式来实现的。
// 这些线程必须通过在彼此之间发送消息，从而向另一个线程请求信息或者将信息提供给另一个线程。
// Go语言中的goroutine和channel机制就是基于CSP模型实现的

type MsgData struct {
	ID  int
	Msg string
}

func TestJZ(t *testing.T) {

	// 先造下数据
	resSlice := make([]*MsgData, 0)
	for i := 0; i < 10; i++ {
		resSlice = append(resSlice, &MsgData{ID: i, Msg: fmt.Sprintf("hello:%d", i)})
	}
	wg := sync.WaitGroup{}
	for pos, v := range resSlice {

		wg.Add(1)
		// 竞争关系 尤其是for循环结构体指针
		// go func(pos int) {
		// 	fmt.Println(pos) // 协程执行打印输出也是无序的 竞争关系 pos也要先复制给临时变量，不然是同一个变量，最后一个
		// 	v.ID = pos + 100 // 直接使用v修改成员属性，只会改最后一个
		// 	wg.Done()
		// }(pos)

		tmpPos := pos
		go func(item *MsgData) {
			fmt.Println(tmpPos)
			item.ID = tmpPos + 100 // 只要结构体指针 也先赋值给临时变量，也不会收到竞争影响
			item.Msg = fmt.Sprintf("jz:%d", item.ID)
			wg.Done()
		}(v)
	}

	wg.Wait()
	for pos, v := range resSlice {
		fmt.Printf("pos:%d, v:%+v\n", pos, v)
	}
}

// pos:0, v:&{ID:0 Msg:hello:0}
// pos:1, v:&{ID:1 Msg:hello:1}
// pos:2, v:&{ID:2 Msg:hello:2}
// pos:3, v:&{ID:3 Msg:hello:3}
// pos:4, v:&{ID:4 Msg:hello:4}
// pos:5, v:&{ID:5 Msg:hello:5}
// pos:6, v:&{ID:6 Msg:hello:6}
// pos:7, v:&{ID:7 Msg:hello:7}
// pos:8, v:&{ID:8 Msg:hello:8}
// pos:9, v:&{ID:106 Msg:hello:9}

// pos:0, v:&{ID:100 Msg:jz:100}
// pos:1, v:&{ID:101 Msg:jz:101}
// pos:2, v:&{ID:102 Msg:jz:102}
// pos:3, v:&{ID:103 Msg:jz:103}
// pos:4, v:&{ID:104 Msg:jz:104}
// pos:5, v:&{ID:105 Msg:jz:105}
// pos:6, v:&{ID:106 Msg:jz:106}
// pos:7, v:&{ID:107 Msg:jz:107}
// pos:8, v:&{ID:108 Msg:jz:108}
// pos:9, v:&{ID:109 Msg:jz:109}

// 使用并发编程提高响应的案例
func processDeal(val int) int {
	time.Sleep(1 * time.Second) // 耗时久一点，协程执行时的并发才能体现
	return val * 2
}

func TestBingfa(t *testing.T) {
	start := time.Now()

	arrSlice := []int{1, 2, 3, 4, 5}
	var result []int
	for pos, v := range arrSlice {
		_ = pos
		newVal := processDeal(v)
		result = append(result, newVal)
	}
	spew.Printf("result:%+v \n", result)

	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}

// result:[2 4 6 8 10]
// 程序执行耗时(s): 5.006926209

func TestBingfa2(t *testing.T) {
	start := time.Now()
	arrSlice := []int{1, 2, 3, 4, 5}
	var result []int
	// 改为等待锁
	wg := sync.WaitGroup{}
	for pos, v := range arrSlice {
		wg.Add(1)
		_ = pos
		go func(val int) {
			defer wg.Done()
			newVal := processDeal(val)
			result = append(result, newVal)
		}(v)
	}
	wg.Wait()
	spew.Printf("result:%+v \n", result)

	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}

// 问题1 竞态条件 并发情况下保护共享资源result切片
// 所以需要 新值传给协程参数 不然就是result:[10 10 10 10 10]
// result:[4] result:[8] 随机了 result:[2 8] 等
// 程序执行耗时(s): 1.000624583

// 单机锁
func TestBingfa3(t *testing.T) {
	start := time.Now()

	arrSlice := []int{1, 2, 3, 4, 5}
	var result []int

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for pos, v := range arrSlice {
		wg.Add(1)
		_ = pos
		// 由于 goroutine 完成顺序不确定，结果顺序随机
		go func(val int) {
			defer wg.Done()
			// mu.Lock() 锁放前和放后的区别
			newVal := processDeal(val)
			mu.Lock() // 只保护共享资源 由于互斥锁保证了append操作的互斥性，所以每个处理结果都会被安全地添加到result切片中。
			result = append(result, newVal)
			mu.Unlock()
		}(v)
	}
	wg.Wait()
	spew.Printf("result:%+v \n", result)

	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}

// result:[4 10 6 2 8]
// 程序执行耗时(s): 1.001229875

// 原子 错误用法 正确的用法
func TestBingfa4(t *testing.T) {
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

// 拆分成非共享内存
func TestBingfa5(t *testing.T) {
	start := time.Now()

	arrSlice := []int{1, 2, 3, 4, 5}
	result := make([]int, len(arrSlice))

	wg := sync.WaitGroup{}
	for pos, v := range arrSlice {
		wg.Add(1)

		// 由于 goroutine 完成顺序不确定，结果顺序随机
		go func(v int, val *int) {
			defer wg.Done()
			newVal := processDeal(v)
			// 直接把引用类型 result拆分为非共享的类型 指针
			*val = newVal

		}(v, &result[pos])

	}
	wg.Wait()
	spew.Printf("result:%+v \n", result)

	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}

// result:[2 4 6 8 10]
// 程序执行耗时(s): 1.001310292

// 使用chan 通过通道实现一把锁
func TestBingfa6(t *testing.T) {
	start := time.Now()

	arrSlice := []int{1, 2, 3, 4, 5}
	result := make([]int, 0)

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, 1) // 默认0无缓冲通道，没有另外协程接收 一放入就会堵塞
	ch <- struct{}{}
	for pos, v := range arrSlice {
		wg.Add(1)
		_ = pos
		// 由于 goroutine 完成顺序不确定，结果顺序随机
		go func(v int) {
			defer wg.Done()
			newVal := processDeal(v)
			<-ch
			result = append(result, newVal)
			ch <- struct{}{}
		}(v)

	}
	wg.Wait()
	spew.Printf("result:%+v \n", result)

	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}

// result:[8 2 10 4 6]
// 程序执行耗时(s): 1.00154425
