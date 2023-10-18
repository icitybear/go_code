package main

import (
	"app_sdk/applog"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// 直接用
func One() {
	user := struct {
		Name  string
		Level int64
	}{Name: "citybear", Level: 2}
	// 使用第三方包 log.Info()
	// {"level":"info","user":{"Name":"chihuo","Level":1},"name":"chihuo","level":1,"time":"2022-10-12T22:23:49+08:00","message":"打印用户信息"}
	// 如果不使用Msg 使用send不会有message
	// log.Info().Interface("user", user).Str("name", user.Name).Int64("level", user.Level).Msg("打印用户信息")
	log.Info().Interface("user", user).Str("name", user.Name).Int64("level", user.Level).Send()
}

// zerolog.New
func Two() {
	// zerolog的第二种用法  输出彩色的日志
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	//
	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend|os.ModePerm)

	if err != nil {
		fmt.Printf("open log file err: %v\n", err)
		return
	}
	// 屏幕输出，文件也有  多输出
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

	// 记录日志带时间 Timestamp  又返回logger
	logger := zerolog.New(multi).With().Timestamp().Logger()

	// {"level":"info","time":"2022-10-12T22:31:49+08:00","message":"Hello World!"}
	// 日志信息记录了字符串
	logger.Info().Msg("Hello World!")
	logger.Debug().Msg("Debug!")
}

func Three() {
	applog.Init("app", "pay")
	applog.Logger("app").Info().Msg("user info.")
	applog.Logger("pay").Info().Msg("pay info.")
}

func main() {
	//One()
	//Two()
	Three()
}
