package panic_recover

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"
)

// 定义错误
var CustomError = errors.New("citybear custom error")

func TestPanicVxExit(t *testing.T) {

	defer func() {
		fmt.Println("defer")
	}()
	fmt.Println("Start")
	//recover 相当于 try catch 捕捉到的异常
	os.Exit(-1) //不会调⽤ defer 指定的函 不输出当前调⽤栈信息
	fmt.Println("End")
}

func TestPanicVxPanic(t *testing.T) {

	defer func() {
		fmt.Println("defer")
		//容易形成僵⼫服务进程，导致 health check 失效
		if err := recover(); err != nil {
			//可以针对具体err 进行处理，建议重启进程 额不是恢复错误 避免僵尸进程
			_, file, line, _ := runtime.Caller(3)
			fmt.Printf("file %v line %d", file, line)
			fmt.Println("recovered from ", err) //err  就是errors
		}
	}()
	fmt.Println("Start")
	//recover 相当于 try catch 捕捉到的异常
	panic(CustomError)
	// 即使recover 捕捉了 但是panic下面的代码还是不会执行
	fmt.Println("End")
}

func TestPanicVxPanic2(t *testing.T) {
	defer fmt.Println("宕机后要做的事情1")
	defer fmt.Println("宕机后要做的事情2")
	panic("宕机")
	defer fmt.Println("宕机后要做的事情3")
}
