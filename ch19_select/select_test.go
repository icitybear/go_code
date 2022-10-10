package select_test

import (
	"fmt"
	"testing"
	"time"
)

func service() string {
	time.Sleep(time.Millisecond * 500) //50 就不超时
	return "Done"
}

func AsyncService() chan string {
	retCh := make(chan string, 1)
	go func() {
		ret := service()
		fmt.Println("returned result.")
		retCh <- ret
		fmt.Println("service exited.")
	}()
	return retCh
}

func TestSelect(t *testing.T) {
	select {
	//不同channel 收到消息执行顺序与case无关 如果都没收到 就默认走default
	case ret := <-AsyncService(): //通过协程返回的是DONE 如果没超时就走着流程
		t.Log(ret)
	case <-time.After(time.Millisecond * 700): //实现超时
		t.Log("time out2")
		//t.Error("time out") //只要1000（xxx）小于500 就会报错超时
	}
}
