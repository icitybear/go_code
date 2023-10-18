package applog

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// 存放 value为指针比较节省空间
var loggers = make(map[string]*zerolog.Logger)

func init() {
	fmt.Print("applog init \n")
}
func Init(names ...string) {
	for _, name := range names {
		// 相对目录是相对于执行文件的路径 比如现在是 main three Init
		openFile, err := os.OpenFile(fmt.Sprintf("./logs/%s.log", name), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend|os.ModePerm)
		if err != nil {
			fmt.Printf("open log file err: %v\n", err)
			return
		}
		//loggers[name] = &zerolog.New(openFile) //这里存的已经是指针了
		tmp := zerolog.New(openFile)
		loggers[name] = &tmp
	}
}

func Logger(name string) *zerolog.Logger {
	instance := loggers[name] //这里存的已经是指针了
	// 若干是+v连结构体字段也有 &{w:{Writer:0xc0000b0020} level:-1 sampler:<nil> context:[] hooks:[] stack:false}
	// v  &{{0xc000014030} -1 <nil> [] [] false}
	fmt.Printf("instance type %T %v\n", instance, instance)
	return instance
}
