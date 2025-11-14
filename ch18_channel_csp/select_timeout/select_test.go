package select_timeout

import (
	"fmt"
	"testing"
	"time"
)

func service() string {
	time.Sleep(time.Millisecond * 1000) // 1000以下就不超时
	return "Done"
}

func AsyncService() chan string {

	retCh := make(chan string, 1)
	go func() {
		//ret := service()
		time.Sleep(time.Millisecond * 500)
		msg := "citybear"

		fmt.Println("returned result.")
		retCh <- msg //由于chan是有容量的 所以不阻塞
		fmt.Println("service exited.")
	}()
	return retCh
}

// 监听 单流程下⼀个go只能监控⼀个channel的状态，select可以完成监控多个channel的状态
func TestSelect(t *testing.T) {
	select {
	//不同channel 收到消息执行顺序与case无关 如果都没收到 就默认走default.如果没有 default 子句,select 将阻塞,直到某个通道可以运行;
	case ret := <-AsyncService(): //通过协程返回的是DONE 如果没超时就走着流程
		t.Log(ret)
	case <-time.After(time.Millisecond * 1000): //实现超时
		t.Log("time out2")
		//t.Error("time out") //只要1000（xxx）小于500 就会报错超时
	}
}

// time.Sleep 改成 time.After 可以降低 cpu 使用率
// time.Sleep 不会让度 cpu 使用权（在这个代码里，cpu 被白白浪费 100ms 的时间）
// time.After 会让度出来

// 监听 time.After 的同时加上 ctx.Done 就可以实现安全退出机制
