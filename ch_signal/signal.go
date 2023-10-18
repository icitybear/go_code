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

func main() {
	// work()       //验证后台进程  再等待信号 输出具体的信号
	// os.Exit(101) //exit status 101

	// 以下优雅退出 defer + 信号监听处理  报错非自然退出 或自然退出 都能调用 myquit
	fmt.Println("main start")

	defer func() {
		fmt.Println("bye main from defer")
		myquit() // 处理退出前  比如服务发现里的 服务注销退出除了
	}()

	sig := make(chan os.Signal)
	// 信号量  ctrl+c SIGINT  服务重新加载状态 SIGHUP
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
					os.Exit(int(i)) // 打印信号量的值
				} else {
					os.Exit(0)
				}
			}
		}
	}()

	// 自然退出时
	wait := make(chan bool)
	go func() {
		for {
			time.Sleep(10000 * time.Millisecond)
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
