package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang-jwt/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

// func main() {
// 	// jwt是包名
// 	jwt.Secret = "123456"

// 	payload := struct {
// 		Id   string `json:"id"`
// 		Name string `json:"name"`
// 	}{Id: "1000", Name: "chihuo"}

// 	sign, err := jwt.Sign(payload)
// 	if err != nil {
// 		fmt.Printf("err %v\n", err)
// 		return
// 	}

// 	fmt.Printf("sign is %s\n", sign)

// 	err = jwt.Verify(sign)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Printf("verify ok")
// }

type UserInfoClaims struct {
	Id                     int      `json:"id"`
	LoginName              string   `json:"loginName"`
	GidList                []string `json:"gidList"`
	jwtv4.RegisteredClaims          // 内嵌标准的声明
}

func main() {

	myClaims := UserInfoClaims{
		1,
		"csx",
		[]string{"a", "b", "c"},
		jwtv4.RegisteredClaims{
			ExpiresAt: jwtv4.NewNumericDate(time.Now().Add(time.Hour * 24)), // 设置JWT过期时间,此处设置为2小时
			IssuedAt:  jwtv4.NewNumericDate(time.Now()),
			Issuer:    "citybear", // 设置签发人
		},
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	// 加盐 []byte
	token1, err := claims.SignedString([]byte("hello"))
	if err != nil {
		fmt.Printf("[AuthJWT]#%v", err)
		return
	}
	fmt.Println(token1)
	tokenString := token1
	token, err := jwt.ParseWithClaims(tokenString, &UserInfoClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return []byte("hello"), nil // 必须使用[]byte
	})
	if err != nil {
		fmt.Println("[ParseWithClaims]#%v", err)
		return
	}

	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*UserInfoClaims); ok && token.Valid { // 校验token

		fmt.Printf("at time:%d\n", claims.IssuedAt.Unix())
		spew.Println(claims)
		return
	}
	fmt.Println("token invalid")
	return
}
