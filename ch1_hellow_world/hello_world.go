package main

import (
	"fmt"
	"os"
)

var golabName string

func init() {
	golabName = "city"
}

func loadconfig() {
	var golabName = "bear"
	// 此时局部优先全局 这里也是局部变量
	fmt.Println(golabName)
}

func loadconfig2() {
	// 未声明局部 所以直接使用全局
	golabName = "HAHAH"
}

func main() {
	fmt.Println(golabName)
	loadconfig()
	fmt.Println(golabName)
	loadconfig2()
	fmt.Println(golabName)

	if len(os.Args) > 1 {
		fmt.Println("Hello World", os.Args[1])
	}
	os.Exit(4)
}
