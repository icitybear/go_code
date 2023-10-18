package main

import "fmt"

// https://segmentfault.com/a/1190000017846997
// https://blog.51cto.com/u_14524391/2891607
// 文件命名约定可以在go build 包里找到详细的说明，简单来说，就是源文件包含后缀：_$GOOS.go，那么这个源文件只会在这个平台下编译
func main() {
	// 按文件后缀名来 window平台会执行 linux会执行
	ret := GetDefaultPath()
	fmt.Print(ret)
}
