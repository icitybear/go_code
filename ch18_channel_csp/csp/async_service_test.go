package concurrency

import (
	"fmt"
	"testing"
	"time"
)

// 同步阻塞channel make(chan Type)
// 带缓冲的阻塞  make(chan Type, int)
// channel <- val 输入数据
// <- channel 接收数据并丢弃
// x := <- channel 接收数据, 并赋值给x
// x, ok := <- channel 接收数据, 并赋值给x，同时检测通道是否为空或者已关闭

// range 通道 相当于替换for的简洁写法

// 无缓冲的同步阻塞
func TestX(t *testing.T) {
	c := make(chan int)
	go func() {
		defer fmt.Println("goroutine stop")
		fmt.Println("goroutine running")
		c <- 66 // 如果外部没有接收 就会同步阻塞 stop是在num赋值后继续执行的
	}()
	fmt.Println("main receive before")
	num := <-c // 同时也同步阻塞在等chan输入数据
	fmt.Println("num = ", num)
}

// 有缓冲的
func TestB(t *testing.T) {
	c := make(chan int, 6)

	go func() {
		defer fmt.Println("goroutine stop")
		fmt.Println("goroutine running")
		for i := 1; i <= 6; i++ {
			fmt.Println("goroutine input=", i)
			c <- i // （容量满3个）开始阻塞 到i=4的时候阻塞，等
		}
		// 关闭后 还能从chan c读取数据 直到读到关闭
		// 而且如果是容量刚好是6 一下子子协程就结束了，不close的话，main也会一直等待，不知道子协程结束了，从nil chan读取会panic
		close(c) // 必须关闭 ok==false , 不然就是 零值, true 一直死循环
	}()

	fmt.Println("main receive before")
	// 死循环 一直获取
	// for {
	// 	// 加上超时可以看出缓冲容量的作用
	// 	time.Sleep(2 * time.Second)
	// 	if num, ok := <-c; ok {
	// 		fmt.Println("num = ", num)
	// 	} else {
	// 		break
	// 	}
	// }
	// 简化成
	for num := range c {
		fmt.Println("num = ", num)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("main receive stop")
}
