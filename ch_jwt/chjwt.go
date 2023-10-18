package main

import (
	"fmt"
	// module下（jwt） 的包jwt2目录层级要一致
	jwt "jwt/jwt2"
	//  github.com/langwan/go-jwt-hs256 第三方组件（库） 也是module名 该层级下有直接xxx.go文件（package包名,与文件名一致）可以直接用
)

func main() {
	// jwt是包名
	jwt.Secret = "123456"

	payload := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{Id: "1000", Name: "chihuo"}

	sign, err := jwt.Sign(payload)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}

	fmt.Printf("sign is %s\n", sign)

	err = jwt.Verify(sign)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("verify ok")
}
