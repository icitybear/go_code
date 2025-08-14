module jwt_test

go 1.19

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang-jwt/jwt/v4 v4.5.2
	jwt v1.0.0
)

replace jwt => ./jwt // 指向本地目录
