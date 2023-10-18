package main

import (
	"fmt"
	//使用sdk下的test_app 只需要sdk一个go.mod
	//自己调用时 go.mod require sdk v0.0.0  replace sdk => ./sdk
	"sdk/test_app"

	//这里也可以使用custom/test_app 只要go.mod那边本地文件路径一致，理解成导入目录
	//然后代码里使用的时候 包名.函数
	// 第三方包 需要先 go get下 然后go mod tidy修改go.sum
	uuid "github.com/satori/go.uuid"
)

func main() {
	fmt.Println("begin")
	test_app.GetCeshi()
	fmt.Println(test_app.AppSize)
	ver := 10
	size, version, name := test_app.GetOut(ver)
	fmt.Println(size, version, name)
	u1 := uuid.NewV4()
	fmt.Println(u1)

	// 返回有几个参数就要几个接收 可以用占空符号
	if c, err := test_app.GetCeshi2(); err != nil {
		fmt.Println(err)
		_ = c
	}
	//_, err := test_app.GetCeshi2()

}
