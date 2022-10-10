package channel_close

import (
	"fmt"
	"sync"
	"testing"
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
