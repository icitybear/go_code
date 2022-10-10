package concurrency

import (
	"fmt"
	"testing"
	"time"
)

func service() string {
	fmt.Println("service is start")
	time.Sleep(time.Millisecond * 50)
	fmt.Println("service is done")
	return "chan data service Done"
}

func otherTask() {
	fmt.Println("otherTask start")
	time.Sleep(time.Millisecond * 100)
	fmt.Println("Task is done")
}

func TestService(t *testing.T) {
	//非异步 按顺序
	fmt.Println(service())
	otherTask()
}

func AsyncService() chan string { //chan通道返回
	retCh := make(chan string, 1) //buffchan定容量的 service exited不用等到外面用通道 就继续往下执行 不等待
	//retCh := make(chan string) //等待client接收通道数据
	go func() { //开启一个协程
		ret := service() //异步执行service了 这个返回要50毫秒 看chan类型是否等待 数据放入通道
		fmt.Println("returned result.")
		retCh <- ret //"Done"结果放到通道 chan的2种方式 如果是一定要有数据ret 那么 service exited. 一定是在<-retCh （输出Done）之后调用
		fmt.Println("service exited.")
	}()
	fmt.Println("AsyncService return")
	return retCh
}

//
func TestAsynService(t *testing.T) {
	retCh := AsyncService() // 通道结果在协程里
	otherTask()
	fmt.Println(<-retCh) //通道有数据 立马打印了
	time.Sleep(time.Second * 1)
}
