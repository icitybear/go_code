package hashtest

import (
	"fmt"
	"testing"
	"time"

	"math/rand"

	"github.com/spaolacci/murmur3"
)

func TestParse(t *testing.T) {
	scale := 50 //概率
	// 演示随机生成1w次hash的分布情况
	hit1 := 0
	hit2 := 0
	uhit1 := 0
	uhit2 := 0
	for i := 0; i < 100; i++ {
		text := randStr() // 生产的随机字符串
		hash := murmur3.Sum32([]byte(text))
		fmt.Printf("str: %s, Hash: %d\n", text, hash)
		// 模是1000的情况
		if hash%1000 < uint32(scale)*10 {
			hit1 += 1
		} else {
			uhit1 += 1
		}
		// 模是100的情况
		if hash%100 < uint32(scale) {
			hit2 += 1
		} else {
			uhit2 += 1
		}
	}

	fmt.Printf("hit1: %d, uhit1: %d \n hit2: %d, uhit2: %d\n", hit1, uhit1, hit2, uhit2)
}

func randStr() string {
	n := 32
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
