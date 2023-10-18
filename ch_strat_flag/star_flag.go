package main

import (
	"flag"
	"fmt"
	"os"
)

type Options struct {
	WorkDir string
	Port    int
	Debug   bool
}

func main() {
	wd, _ := os.Getwd()
	// 通过启动参数配置 目录 环境 端口
	opt := Options{WorkDir: wd, Port: 80, Debug: false}

	//输入bool
	flag.BoolVar(&opt.Debug, "debug", false, "true is debug")
	// 参数类型化    参数名port 默认值80 提示语web server port
	flag.IntVar(&opt.Port, "port", 80, "web server port")
	flag.StringVar(&opt.WorkDir, "dir", wd, "work directory")

	flag.Parse() //这行必须要 解析参数

	// 设置参数

	fmt.Println(opt)

}
