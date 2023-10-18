package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var mysyncmap sync.Map //map时非线程安全，并发要使用线程安全的sync.Map

func main() {
	//gin框架
	g := gin.Default()
	//gin的网络请求路由 post请求方式
	g.POST("/login", func(context *gin.Context) {
		//每次请求login都生成uuid
		token := uuid.NewV4().String()
		//保存 服务程序执行的时候，希望mysyncmap是线程安全的，数据都是独立，保存了
		mysyncmap.Store(token, token)
		// context上下文 匿名结构体的数据
		context.JSON(200, struct {
			Token string `json:"token"`
		}{Token: token})
	})
	g.POST("/doing", func(context *gin.Context) {

		name := context.PostForm("name")
		fmt.Print(name)
		req := struct {
			Token string `json:"token"`
		}{}
		// context 绑定了参数  这里用的是post json数据格式
		context.ShouldBind(&req)
		fmt.Print(req)
		//查询Load
		if _, ok := mysyncmap.Load(req.Token); ok {
			//也可以String输出
			context.JSON(200, struct {
				Message string `json:"message"`
			}{Message: "find"})
		} else {
			context.JSON(403, struct {
				Message string `json:"message"`
			}{Message: "not find"})
		}
	})
	// 404 page not found
	// 绑定再8080端口
	g.Run(":8080")
}
