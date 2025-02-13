package main

import (
	"fmt"
	"os"
	"os/signal" //信号量
	"syscall"   //系统包
	"time"
)

func work() {
	//信号量用途，监听信号量 处理完自己的进程收尾..监听一些别的信号量，重新定义等 10 12
	// 比如日志分割用到，优雅退出用到  经常监听USR1 (30) USR2 (31)  ctrl+c  signal:interrupt
	// kill -USR1 12345
	fmt.Printf("start\n")
	pid := os.Getpid() //当前PID
	fmt.Println(pid)
	sig := make(chan os.Signal) //监听信号量的通道 收到信号通道就有值
	die := make(chan bool)      // 用来阻塞
	// 监听信号量signal.Notify
	signal.Notify(sig, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		for {
			s := <-sig
			fmt.Printf("recv sig %d\n", s) //输出信号量的值
			if s == syscall.SIGUSR1 {
				die <- true //为USR1时通知外面main协程
			}
		}
	}()
	<-die
	fmt.Printf("exit\n")
}

// 如果程序是被手动终止（如 kill -9），可以捕获信号并实现优雅退出。注意，SIGKILL（kill -9）无法被捕获，但可以捕获其他信号（如 SIGTERM）。
// kill -USR1 12345 USR1 (30) USR2 (31)是自定义信号量 可以用来处理 例如触发日志轮转、重新加载配置或执行其他自定义操作。
// kill -15 terminated （SIGTERM）
// ctrl+c =》2  SIGINT

// 如果你或其他人手动使用 kill -9 <PID> 终止了程序，也会导致 signal: killed
// 信号处理函数的幂等性：
//	确保信号处理函数是幂等的，即多次调用不会产生副作用。
// 信号处理的阻塞：
//	信号处理函数应尽量快速完成，避免阻塞主程序。如果需要执行耗时操作，可以将其放到单独的goroutine中执行。
// 跨平台兼容性：
//	SIGUSR1和SIGUSR2是Unix/Linux特有的信号，在Windows上不可用。如果需要在Windows上运行，请使用其他机制（如HTTP接口或文件监听）来实现类似功能

func main() {
	// work()       //验证后台进程  再等待信号 输出具体的信号
	// os.Exit(101) //exit status 101

	// 以下优雅退出 defer + 信号监听处理  报错非自然退出 或自然退出 都能调用 myquit
	fmt.Println("main start running. PID:", os.Getpid())
	defer func() {
		fmt.Println("bye main from defer")
		myquit() // 处理退出前  比如服务发现里的 服务注销退出除了
	}()

	sig := make(chan os.Signal)
	// 信号量  ctrl+c =》2  SIGINT  服务重新加载状态 SIGHUP
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		// 用range去取信号量通道的值 监听信号 阻塞且循环 		for { s := <-sig .....}
		for s := range sig {
			switch s {
			// 非自然退出时
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP:
				// 调用自己退出处理函数
				myquit()
				if i, ok := s.(syscall.Signal); ok {
					// os.Exit(int(i)) // 打印信号量的值
					// 如果只是单纯打印 并没有os.Exit 还是会执行外面的 defer
					fmt.Println(int(i))
				} else {
					os.Exit(0)
				}
			}
		}
	}()

	defer func() {
		fmt.Println("tttttt")
	}()

	// 自然退出时
	wait := make(chan bool)
	go func() {
		for {
			time.Sleep(20000 * time.Millisecond)
			//close(wait)
			wait <- true
		}
	}()
	<-wait

	fmt.Println("main end")
}

func myquit() {
	fmt.Println("\n myquit 成功退出")
}

// 特性	kill -15（SIGTERM）	kill -9（SIGKILL）
// 信号类型	优雅终止信号	强制终止信号
// 进程能否捕获	可以捕获并执行清理操作	无法捕获或忽略
// 默认行为	终止进程	立即终止进程
// 使用场景	希望进程优雅退出时使用	进程无响应或无法通过 SIGTERM 终止时使用
// 资源清理	进程可以执行清理操作	进程无法执行清理操作
