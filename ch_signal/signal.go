package main

import (
	"context"
	"fmt"
	"os"
	"os/signal" //信号量
	"sync"
	"syscall" //系统包
	"time"

	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"
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

func main1() {
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

	// tag: 如果是子协程里 就要使用单独ctx
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

func childProcess(ctx context.Context) {

	// gCtx, cancel := context.WithCancel(ctx)
	// defer cancel()
	// // 启用信号量监听中断 先处理自定义事件
	// sig := make(chan os.Signal)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
	// go func(ctx context.Context) {
	// 	// 用range去取信号量通道的值 监听信号 阻塞且循环 		for { s := <-sig .....}
	// 	for s := range sig {
	// 		// 每次取信号量前判断下程序是否结束了
	// 		select {
	// 		case <-ctx.Done(): // 接收取消通知 如果超时也会收到 ctx.Err() context deadline exceeded
	// 			return // 结束
	// 		default:
	// 			switch s {
	// 			// 非自然退出时
	// 			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP:
	// 				// 调用自己退出处理函数
	// 				if i, ok := s.(syscall.Signal); ok {
	// 					// 如果只是单纯打印 并没有os.Exit 还是会执行外面的 defer
	// 					fmt.Printf("[ExecuteRecord] 程序非自然退出 捕捉信号量#%d, time:%+v", int(i), time.Now().Format("2016-01-02 15:04:05"))
	// 					return
	// 				}
	// 			}
	// 		}
	// 	}
	// }(gCtx)

	// 继续开启子协程 处理信号量
	bT := time.Now()
	fmt.Printf("[childProcess]: Run start:%v \n", bT)

	list := map[string]int{
		"sun1": 1,
		"sun2": 2,
		"sun3": 3,
		"sun4": 4,
		"sun5": 5,
		"sun6": 6,
	}

	limit := rate.Every(time.Second / time.Duration(1))
	limiter := rate.NewLimiter(limit, 1)

	wg := sync.WaitGroup{}
	concurrent := semaphore.NewWeighted(3) // 限制最大并发数 最大资源数3
	// ctx上层被cancel 但是没被监听 照旧执行
	for k, v := range list {
		concurrent.Acquire(ctx, 1)
		wg.Add(1)
		_ = v
		go sunProcess(ctx, k, &wg, limiter, concurrent)
	}

	wg.Wait()
	eT := time.Since(bT)
	fmt.Printf("[childProcess]: Run time:%v \n", eT)
}

func sunProcess(ctx context.Context, k string, wg *sync.WaitGroup, limiter *rate.Limiter, concurrent *semaphore.Weighted) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[sunProcess]#Runtime panic caught: %v\n", err) // 捕捉recover
		}
	}()

	// defer concurrent.Release(1) // n 个资源可用了 回到总资源池 panic的情况时 s.cur < 0 当前已被使用的资源
	// 如果ctx最顶层被cancel了 优化版
	defer func() {
		// 阻塞等待
		select {
		case <-ctx.Done():
			fmt.Printf("[sunProcess]#%s ctx done:%v \n", k, ctx.Err())
		default:
			fmt.Printf("[sunProcess]#%s success \n", k)
			concurrent.Release(1)
		}
	}()

	defer wg.Done()

	num := 0
	for {
		// 抢等待拿到令牌
		err := limiter.WaitN(ctx, 1)
		if err != nil {
			break
		}
		fmt.Printf("[sunProcess]#%s 抢到令牌 time:%s \n", k, time.Now().Format("2006-01-02 15:04:05"))
		num = num + 1
		if num > 2 {
			fmt.Printf("[sunProcess]#%s 退出 time:%s \n", k, time.Now().Format("2006-01-02 15:04:05"))
			break // 执行10次后就退出
		}
	}

}

func main() {
	// 模拟收到信号退出 和超时退出
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 启用信号量监听中断 当前主协程
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	gCtx, _ := context.WithCancel(ctx) //_ gCancel
	go func(gctx context.Context) {
		// 用range去取信号量通道的值 监听信号 阻塞且循环 		for { s := <-sig .....}
		for s := range sig {
			// 每次取信号量前判断下程序是否结束了
			select {
			// 这里的子ctx
			case <-gctx.Done(): // 接收取消通知 如果超时也会收到 context deadline exceeded
				fmt.Printf("[main] gctx done#%v, time:%+v \n", gctx.Err(), time.Now().Format("2016-01-02 15:04:05"))
				return // 结束
			default:
				switch s {
				// 非自然退出时
				case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP:
					// 调用自己退出处理函数
					if i, ok := s.(syscall.Signal); ok {
						// 如果只是单纯打印 并没有os.Exit 还是会执行外面的 defer和程序
						fmt.Printf("[main] 程序非自然退出 捕捉信号量#%d, time:%+v \n", int(i), time.Now().Format("2016-01-02 15:04:05"))
						// os.Exit(1)
						cancel() // 收到信号 取消main的ctx 子协程报错panic [sunProcess]#Runtime panic caught: semaphore: released more than held
						// gCancel() // 通过取消gctx的就不会有这种情况
						return
					}
				}
			}
		}
	}(gCtx)

	go childProcess(ctx) // 类似

	// 阻塞等待
	select {
	case <-ctx.Done():
		fmt.Println("main ctx done", ctx.Err())
	}
	time.Sleep(2 * time.Second) // 等待子协程打印完消息
}
