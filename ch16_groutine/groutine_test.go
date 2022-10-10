package groutine_test

import (
	"fmt"
	"testing"
	"time"
)

func TestGroutine(t *testing.T) {

	for i := 0; i < 20; i++ {

		// i 值传递 协程里的变量地址是不一样的 无竞争关系 就不会有共享内存了
		go func(i int) {
			//time.Sleep(time.Second * 1) //外面程序执行时间只要大于 协程，就会有输出了
			fmt.Println(i)
		}(i)

		//竞争关系 运行时间快 导致输出的全是20 共享内存了
		// go func() {
		// 	fmt.Println(i)
		// }()
	}

	time.Sleep(time.Millisecond * 100) //加这个是因为 外面执行完的速度 超过了协程 这样就没输出了
}
