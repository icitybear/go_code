package jwt2

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// 只支持了HS256 散列算法
const (
	HS256 = "HS256"
)

var alg = HS256

var Secret string

func hs256(secret, data []byte) (ret string, err error) {
	// hmac库带上secret 初始化一个散列
	hasher := hmac.New(sha256.New, secret)
	_, err = hasher.Write(data) //填进数据
	if err != nil {
		return "", err
	}
	r := hasher.Sum(nil) //得到一个hash散列值
	// 再base64化
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEwMDAiLCJuYW1lIjoiY2hpaHVvIn0.1egKEEN3IRaK1wblqGzkJQ5wciKrssslqAAiLXo8iTA
	return base64.RawURLEncoding.EncodeToString(r), nil
}

// 签名
func Sign(payload interface{}) (ret string, err error) {
	h := header{
		Alg: alg,
		Typ: "JWT",
	}
	marshal, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	// header头部分 base64
	bh := base64.RawURLEncoding.EncodeToString(marshal)

	marshal, err = json.Marshal(payload)
	if err != nil {
		return "", err
	}
	// payload头部分 base64
	bp := base64.RawURLEncoding.EncodeToString(marshal)
	//2部分base64 . 连接后 hs256
	s := fmt.Sprintf("%s.%s", bh, bp) //与官网一致

	//Secret 官网支持多中算法 目前我们采取hs256
	ret, err = hs256([]byte(Secret), []byte(s))
	if err != nil {
		return "", err
	}
	// 返回签名后的密文
	return fmt.Sprintf("%s.%s.%s", bh, bp, ret), nil
}

// 校验
func Verify(token string) (err error) {
	parts := strings.Split(token, ".")
	// 前2部分
	data := strings.Join(parts[0:2], ".")
	// 填进数据
	hasher := hmac.New(sha256.New, []byte(Secret))
	_, err = hasher.Write([]byte(data))
	if err != nil {
		return err
	}

	// 第三部分的签名
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return err
	}
	// 同样hasher.Sum(nil) 获取hash散列值  然后跟 传进来的第三部分签名做一个比较
	if hmac.Equal(sig, hasher.Sum(nil)) {
		return nil
	}
	return errors.New("verify is invalid")
}
