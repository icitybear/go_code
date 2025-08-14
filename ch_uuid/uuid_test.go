package ch_uuid

import (
	// 第三方包 需要先 go get下 然后go mod dity修改go.sum
	"fmt"
	"math/rand"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	uuid2 "github.com/google/uuid"

	"github.com/rs/xid"

	"kid" // 自身kid模块包 测试本地模块包 需要引入module require replace

	jwt "kid/jwts" // 使用模块下的某个包
)

// 直接终端使用命令行运行 终端进入ch_uuid目录  如果追加目录
// $ /usr/local/go/bin/go test -timeout 30s -run ^TestKid$
// 8ugpwzvu4p8gsh
// PASS
// ok      ch_uuid 0.021s
func TestKid(t *testing.T) {

	uuid := kid.New().String()
	fmt.Println(uuid)
}

// $ /usr/local/go/bin/go test -timeout 30s -run ^TestJwt$
// sign is eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEwMDAiLCJuYW1lIjoiY2hpaHVvIn0.1egKEEN3IRaK1wblqGzkJQ5wciKrssslqAAiLXo8iTA
// verify okPASS
// ok      ch_uuid 0.022s
func TestJwt(t *testing.T) {
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

// 因为vscode直接打开的code 直接使用ide工具运行测试
// tag: /usr/local/go/bin/go test -timeout 30s -run ^TestKy$ code/ch_uuid
// 会提示 package code/ch_uuid is not in GOROOT (/usr/local/go/src/code/ch_uuid)
//
//		package记得按目录命名
//	 /usr/local/go/bin/go test -timeout 30s -run ^TestKy$ 才准确
func TestKy(t *testing.T) {

	// 生成唯一ID
	id := xid.New()
	// 将唯一ID转换为字符串
	idString := id.String()
	// 打印唯一ID
	fmt.Println(idString) // 22位 cnpverd315ok504649qg

	u1 := uuid.NewV4()
	fmt.Println(u1) // 36 2783b3c3-cd17-4bb9-9812-b7fefa7433dc

	v, _ := uuid2.New().Value()
	fmt.Println(v) //36 371e8993-341e-4df3-9cc6-f3fc4db90dc4

	uuidStr2 := uuid2.NewString()
	fmt.Println(uuidStr2) // 36 8e7acb31-6670-476b-9564-dc470634d83d

	uuid := uuid2.New().ID()
	fmt.Println(uuid) //10位 2667746093

	n := 32
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	fmt.Println(string(b))

}
