package main

// 第一个go程序，全局变量与局部变量，以及init顺序
import (
	"fmt"
	"os"
)

// 局部变量与全局变量的优先级 全局跨包变量首字母要大写
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
	fmt.Println(golabName)
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
