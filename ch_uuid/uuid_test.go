package chuuid_test

import (
	// 第三方包 需要先 go get下 然后go mod dity修改go.sum
	"fmt"
	"math/rand"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	uuid2 "github.com/google/uuid"

	"github.com/rs/xid"
)

func TestXxx(t *testing.T) {

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
