package test_channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// buff chan定容量的 service exited不用等到外面用通道(接收) 就继续往下执行 不等待
var ch = make(chan int) // 没定容量的 等待client接收通道数据
var sum = 0

// 等待锁
func worker(wg *sync.WaitGroup) {
	// 常见结构
	for {
		// slect多路选择与超时控制
		select {
		//只需要case  接收channel数据  如果没有default则一直等待 (因为没收到默认走default)
		case num, ok := <-ch: //非阻塞式
			fmt.Printf("num is %d \n", num)
			if !ok {
				// 管道结束
				wg.Done()
				return
			}
			sum = sum + num
			//case	<-time.After(time.Millisecond * 100):  //利用时间包  从管道获取
			// wg.Done()
			// return

		}
	}
}

func producer() {
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch) // 通知管道关闭了  使用了<-ch的 num, ok := <-ch  0，false
	// 向关闭的 channel 发送数据，会导致 panic
	// v, ok <-ch; ok 为 bool 值，true 表示正常接受，false 表示通道关闭, 此时v为零值
	// 所有的 channel 接收者都会在 channel 关闭时，⽴刻从阻塞等待中返回且上述 ok 值为 false。
	// 这个⼴播机制常被利⽤，进⾏向多个订阅者同时发送信号。

}

// 使用等待锁的情况下
func TestChan(t *testing.T) {
	// 等待锁  sync.WaitGroup, 因为worker()需要传的时指针
	// wg := sync.WaitGroup{} // {} 正常实例化 非返回指针
	var wg sync.WaitGroup //var 也是常见的结构体 实例化
	fmt.Printf("%T \n", wg)
	wg.Add(1)

	// wg要&取值 表示结构体的指针类型
	// worder不停的等待 chan的输出 直到手动关闭了通知（收信号 收到结束信号） 便结束等待锁
	go worker(&wg)

	// 想象成 producer虫洞 不停给chan 扔完数据并且扔完才关chan (传输信号)
	go producer()
	wg.Wait()
	fmt.Printf("wait sum is %d", sum)
}

// -------------------------------------------------------------------------------//

// 作用把in管道输出的赋值给管道out
func worker1(id int, in chan bool, out chan bool) {
	fmt.Printf("worker1 %d start\n", id)
	<-in
	fmt.Printf("worker1 %d quit\n", id)
	out <- true
}

// chan 读写chan<-只写<-chan 只读
func worker2(id int, in <-chan bool, out chan<- bool) {
	fmt.Printf("worker2 %d start\n", id)
	<-in
	fmt.Printf("worker2 %d quit\n", id)
	out <- true
}

func worker3(id int, in <-chan bool, out chan<- bool) {
	fmt.Printf("worker3 %d start\n", id)
	for {
		fmt.Printf("worker3 %d doing\n", id)
		time.Sleep(200 * time.Millisecond) //不断循环的情况下sleep
		select {
		// 管道有输入 才执行 否则就是default
		case <-in: //bool, ok := <-ch: 收到信号后的操作
			out <- true // 想成一个关闭协程外面的用out chan 等待的程序的信号
			return
		default:
			fmt.Printf("worker3 %d  wait in chan default\n", id)
		}
	}
}

func TestChan2(t *testing.T) {
	// 相当于传输信号了
	in := make(chan bool)
	out := make(chan bool)
	// 开多个协程的情况
	workers := 3

	for i := 0; i < workers; i++ {
		//go worker1(i, in, out)
		//go worker2(i, in, out)
		go worker3(i, in, out) //传值完后 协程的执行顺序是无序的
	}

	// 单独一个协程 不断去往in chan扔输出  相当于一个不断有浏览点击的
	go func() {
		time.Sleep(time.Second)
		close(in) //如果不知道协程数量 直接close  通知所有使用 <-in 等待阻塞的协程 想成一个关闭协程的信号
		// for i := 0; i < workers; i++ {
		// 	in <- true //相当于一个个协程去通知 信号
		// }
	}()

	// 计数器 for 配合chan阻塞 数量技术
	count := 0
	for count < workers {
		// 不断从 out chan 获取数据 用来计数 （协程都执行 执行的把in赋值给out完毕了）
		<-out //阻塞在等着
		count++
	}

	fmt.Println("ok")
}
