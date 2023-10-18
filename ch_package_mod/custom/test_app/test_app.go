package test_app

import "fmt"

// 相当于public protected privalige
var appVersion int
var appName string
var AppSize int

func init() {
	appVersion = 0
	appName = "init"
	AppSize = 0
	fmt.Println(appVersion, appName, AppSize)
}

func GetOut(ver int) (int, int, string) {
	AppSize = 100
	name := "citybear"
	appName = name

	setVersion(ver)
	return AppSize, appVersion, appName
}

func setVersion(ver int) {
	appVersion = ver + 1000
}

func GetCeshi() {
	fmt.Println("ceshi")
}
