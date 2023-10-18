package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256" //使用sha256
)

//微服务 api接口调用 签名与校验。 存在双方互相持有2对公钥密钥 用来互相校验
//或者加密了aes的密钥 组合使用 加密与解密
//rsa非对称加密，解密慢
//私钥加密，公钥解密 就是签名与校验
//公钥加密，私钥解密 就是加密与解密

// 加密 一般工作中难的是 各种格式的文件导入，转成一个的公钥和私钥
func Encrypt(key *rsa.PublicKey, src []byte) (data []byte, err error) {

	h := sha256.New() //crypto/sha256 crypto/md5
	//加密方法 OAEP  rsa加密和解密 PKCS1最早的规范版本v15 大量项目还在使用，虽然有缺陷
	ciphertext, err := rsa.EncryptOAEP(h, rand.Reader, key, src, nil)
	//rsa.EncryptPKCS1v15()
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// 解密
func Decrypt(key *rsa.PrivateKey, src []byte) (data []byte, err error) {
	h := sha256.New()
	// OAEP
	oaep, err := rsa.DecryptOAEP(h, rand.Reader, key, src, nil)
	if err != nil {
		return nil, err
	}
	return oaep, nil
}

// 签名
func Sign(key *rsa.PrivateKey, src []byte) (sign []byte, err error) {
	// 指定算法 也可以是md5 sha1  crypto.MD5 crypto.SHA1
	h := crypto.SHA256
	hn := h.New() //hash自己生成个新的
	hn.Write(src)
	sum := hn.Sum(nil) //得到校验和
	// 签名方法 OAEP  rsa加密和解密 PKCS1最早的规范版本v15 大量项目还在使用，虽然有缺陷
	// rsa.SignPKCS1v15()
	return rsa.SignPSS(rand.Reader, key, h, sum, nil)
}

// 校验
func Verify(key *rsa.PublicKey, sign, src []byte) (err error) {
	h := crypto.SHA256
	hn := h.New() //hash自己生成个新的 不用传
	hn.Write(src)
	sum := hn.Sum(nil)
	// 只要公钥 签名后的密文 是对应的私钥签名的 就没问题
	return rsa.VerifyPSS(key, h, sum, sign, nil)
}
