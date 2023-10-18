package channel_close

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func dataProducer(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i //数据放入chan
		}
		close(ch) //关闭时 v, ok <-ch;从阻塞等待中返回且上述 ok 值为 false。

		wg.Done()
	}()

}

func dataReceiver(ch chan int, wg *sync.WaitGroup) {
	go func() {
		//死循环 直到通道关闭
		for {
			// 通道未关闭 有数据
			if data, ok := <-ch; ok {
				fmt.Println(data)
			} else {
				break
			}
		}

		wg.Done() //跳出循环才完成
	}()

}

// channel 只能一个输入方关闭, 但是能有多个读取，chan是并发安全的
func TestCloseChannel(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int) //必须等待接收
	wg.Add(1)
	dataProducer(ch, &wg) //等待锁 指针传入

	//可以有多个接收者
	wg.Add(1)
	dataReceiver(ch, &wg)
	wg.Add(1)
	dataReceiver(ch, &wg)

	wg.Wait() //等到里面done

}

func TestCloseChannel2(t *testing.T) {
	ch := make(chan int, 3)
	// 发送方
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("发送方: 发送数据 %v...\n", i)
			ch <- i
		}
		fmt.Println("发送方: 关闭通道...")
		close(ch)
	}()
	// 接收方 收到其他协程关闭通道的通知
	for {
		// 加上超时可以看出缓冲容量的作用
		time.Sleep(2 * time.Second)

		num, ok := <-ch // 所有接收通道的地方 都会收到
		if !ok {
			fmt.Println("接收方: 通道已关闭")
			break
		}
		fmt.Printf("接收方: 接收数据: %v\n", num)
	}
	fmt.Println("程序退出")
}
