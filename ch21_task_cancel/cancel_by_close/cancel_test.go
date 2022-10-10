package concurrency

import (
	"fmt"
	"testing"
	"time"
)

func isCancelled(cancelChan chan struct{}) bool {
	//多路选择
	select {
	case v, ok := <-cancelChan: //关闭时 也有数据 {},false
		// 发消息时 因为返回true 然后跳出循环
		fmt.Println(v, ok, "msg")
		return true
	default:
		return false
	}
}

func cancel_1(cancelChan chan struct{}) {
	cancelChan <- struct{}{} //发消息{}
}

func cancel_2(cancelChan chan struct{}) {
	close(cancelChan) //关闭广播 {}
}

func TestCancel(t *testing.T) {
	cancelChan := make(chan struct{}, 1) //等待chan接收

	for i := 0; i < 5; i++ {
		//5个协程 channel等待中
		go func(i int, cancelCh chan struct{}) {
			//循环中 死循环呢
			for {
				if isCancelled(cancelCh) { //等到有数据就跳出循环
					break
				}
				time.Sleep(time.Millisecond * 5)
			}
			fmt.Println(i, "Cancelled")
		}(i, cancelChan)
	}
	cancel_1(cancelChan) //发消息 只有一个chan收到 其他都走select了default
	//cancel_2(cancelChan) //close 所有chann都收到

	time.Sleep(time.Second * 1)
}
