package ch_uuid

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	// 第三方包 需要先 go get下 然后go mod dity修改go.sum
	uuid "github.com/satori/go.uuid"

	uuid2 "github.com/google/uuid"

	"github.com/rs/xid"

	"kid" // 自身kid模块包 测试本地模块包 需要引入module require replace

	jwt "kid/jwts" // 使用模块下的某个包

	"kid/aes"

	"crypto/rand"
	"crypto/rsa"
	appRsa "kid/rsa"
	mrand "math/rand"
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
	mrand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[mrand.Intn(len(charset))]
	}
	fmt.Println(string(b))

}

// $ /usr/local/go/bin/go test -timeout 30s -run ^TestAes$
// str is 0citybear1
// PASS
// ok      ch_uuid 0.019s
func TestAes(t *testing.T) {
	// tag: 双方加密解密(go,java,php) aes算法要确定 1填充方式 2 iv向量的传递 放的位置还是额外传参 3加密的模式 CBC
	//aes128 aes256 强度不一样
	key := "01234567890123456789012345678912" //32位 钥匙
	//待加密的内容 长度随意
	str := "0citybear1"
	// 加密
	encrypt, err := aes.Encrypt([]byte(key), []byte(str))
	if err != nil {
		fmt.Println(err)
		return
	}
	// 解密
	decrypt, err := aes.Decrypt([]byte(key), encrypt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("str is %s\n", string(decrypt))
}

// $ /usr/local/go/bin/go test -timeout 30s -run ^TestRsa$
func TestRsa(t *testing.T) {
	// crypto/rsa 包生成公钥和私钥  公钥密钥长度 2048  4096更长（金融用）
	priKey, err := rsa.GenerateKey(rand.Reader, 4096) //现在是1行，实际工作对接时，获取公钥私钥的解析比较复杂

	if err != nil {
		return
	}
	// 比如加密大文件 先用aes加密，然后aes的key 使用rsa加密保护 存到数据库，取的时候再用rsa解密 （云存储）
	pubKey := priKey.PublicKey // 公钥

	src := "01ddddcccc89"

	// 这里的公钥都是取地址 私钥直接赋值 大小问题
	encrypt, err := appRsa.Encrypt(&pubKey, []byte(src))
	if err != nil {
		fmt.Printf("encrypt err : %v", err)
		return
	}

	decrypt, err := appRsa.Decrypt(priKey, encrypt)
	if err != nil {
		fmt.Printf("decrypt err : %v", err)
		return
	}
	fmt.Printf("解密后的数据: %s\n", decrypt)
	// 90%还是签名和校验
	// 用私钥签名
	sign, err := appRsa.Sign(priKey, []byte(src))
	if err != nil {
		fmt.Printf("sign err : %v", err)
		return
	}

	fmt.Printf("签名后为: %s\n", hex.EncodeToString(sign))

	// 用公钥和签名后的密文  去看校验是否成功
	err = appRsa.Verify(&pubKey, sign, []byte(src))
	if err != nil {
		fmt.Printf("verify err : %v", err)
		return
	} else {
		fmt.Println("校验签名成功")
	}

}

// 解密后的数据: 01ddddcccc89
// 签名后为: 95267041ffa5f04f4faf4ca82e7cf5a07fb1c6253ca9b96fa1da49109d43e62b503edfe15d9f825e5a52a2ab9bf026cb54f8ce83e134dbb4904b7db6c84423e40b4e9764c87ec9c77dcdce8f899e3e0fd2b4854536422439c6e322ef2f8f0bd48fe81c07f9fc9cca3d6ca34bd5ec90e7c9a2523e4244297b1ba2a534a4afbbfefa391342fc5ae9feea14fddca7c35c586efd3ee55a4de025f27657bd6a7ece8102af8e38d1cd98893bba9d03aace92b4ef62d47f3858835d10dec182d07fb84cb38a7328d6a365f9c8e713c300d6d63a3a1b5b044b3d32a60abd4a41c4c6ff07ccdf653784a3dd621731ba91cd798b4010b27ee660c7f63328000ffd1da3bc3be9fc7463259b28984f22a924a92160ba35ef213fa99b27586f0f1f1b062e97666c9146df1714eb71d049a90476efe4642fc96ace95f7f7b393575cc8b5e7b2419813337f0b353e6c45ea56fcd555920cb21a8badbc6ddb90f21c435eb7f5ef78cb05539293a17145f24c765657088dcb60620a2187cf1956663c828944c2e05b11025738c84aef510122d7a6a05ff47832e7b8655b6569a188d309263546ac04df89fa2e1cf734a43803d2ebb11bd7ca743b4273890e1f25797cb6c2886cb1374c6c24d99022e76a39265f2490540660920cb900aa126886ebc582fe92eb838156be2d8e983fd93de8b9cfcdc06daf2bc89a8de99789474c45eb4f60e383bfd7
// 校验签名成功
// PASS
// ok      ch_uuid 2.243s
