package aes

//对称加密aes 速度快
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// 填充方法
func pkcs7Pad(src []byte, blockSize int) (dest []byte, err error) {
	if blockSize <= 0 {
		return nil, errors.New("block size is 0")
	} else if src == nil || len(src) == 0 {
		return nil, errors.New("src is nil")
	}
	// 块有多大 取余数 看补多少
	n := blockSize - (len(src) % blockSize)
	pb := make([]byte, len(src)+n)
	copy(pb, src) // 搬原来的数据 到pb
	// 然后补 n 个缺的
	copy(pb[len(src):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// 去掉填充的方法
func pkcs7Unpad(src []byte, blockSize int) (dest []byte, err error) {

	if blockSize <= 0 {
		return nil, errors.New("block size is 0")
	} else if len(src)%blockSize != 0 {
		return nil, errors.New("src length error")
	} else if src == nil || len(src) == 0 {
		return nil, errors.New("src is nil")
	}

	c := src[len(src)-1] //尾巴的长度

	padLength := int(c)

	if padLength == 0 || padLength > len(src) {
		return nil, errors.New("pad length error")
	}

	for i := 0; i < padLength; i++ {
		if src[len(src)-padLength+i] != c {
			return nil, errors.New("pad content error")
		}
	}

	// 砍掉尾巴
	return src[:len(src)-padLength], nil

}

// aes加密方法
func Encrypt(key, src []byte) (data []byte, err error) {

	block, err := aes.NewCipher(key) //生成的块

	if err != nil {
		return nil, err
	} else if len(src) == 0 {
		return nil, errors.New("src is empty")
	}

	// 工程 流水线 肉 切成标准块  最后一块的时候 对不齐 用这种填充方式凑够 pkcs7 pkcs5填充方式标准 （交接的时候要确定）
	plaintext, err := pkcs7Pad(src, block.BlockSize())

	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// iv向量 保证加密解密更难破解  这里时把iv向量放生成加密块的前面
	iv := ciphertext[:aes.BlockSize]
	// rand.Reader随机生成 iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	// 块与块的加密方式时CBC
	bm := cipher.NewCBCEncrypter(block, iv)
	// 加解密的实质，把处理后的肉块plaintext
	// 填充到缓冲buff里（预留好的数据块  aes块大小（iv向量大小）+密文大小）
	// 填充到iv向量之后aes.BlockSize:
	bm.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// 得到 前面有iv向量 + 加密后的密文 结果
	return ciphertext, nil
}

func Decrypt(key, src []byte) (data []byte, err error) {

	if len(src) < aes.BlockSize {
		return nil, errors.New("data length error")
	}

	// 先找出iv向量  之前放再前面aes块大小
	iv := src[:aes.BlockSize]
	// 后面是 加密后的数据
	ciphertext := src[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key) //生成的块
	if err != nil {
		return nil, err
	}
	// 一样的cbc方式
	bm := cipher.NewCBCDecrypter(block, iv)
	//ciphertext机密的数据 放进去处理后 变成解密了的数据ciphertext
	bm.CryptBlocks(ciphertext, ciphertext)
	// 解密后 结尾多余的那块填充的要去掉
	ciphertext, err = pkcs7Unpad(ciphertext, aes.BlockSize)

	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}
