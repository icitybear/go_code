package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	appRsa "rsa/rsa"
)

func main() {
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
