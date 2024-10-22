package hashtest

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"testing"
	"time"

	"math/rand"

	"github.com/flike/kingshard/core/hack"
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

func TestMd5(t *testing.T) {
	// be19fe18d4363ccfcff17f3d81b649ed caid直接加密后 d4f932cebf2138abd48d2a07f5996a20
	// 42ACABF7-1345-4A7E-A658-E108D8B20B1A
	str := "be19fe18d4363ccfcff17f3d81b649ed"
	h := md5.New()
	h.Write([]byte(str))
	m := hex.EncodeToString(h.Sum(nil))
	println(m)

	h1 := md5.New()
	h1.Write([]byte(strings.ToUpper(str))) // 转大写后加密值就不一样了 有区分大小写
	m1 := hex.EncodeToString(h1.Sum(nil))
	println(m1)
}

func TestHash32(t *testing.T) {
	// system = 'mp-admp-data-job' and message like '%推送基础数据给大数据成功%' and message like '%bhciijcjbcjebici%'
	str := "BA73F597-AC70-4275-BED2-E366CF5F0235" // token
	hash := crc32.ChecksumIEEE([]byte(str)) >> 16 & 0x7FFF
	num := int(hash % 32)
	println(num)
}

// 使用CRC32来计算分表数
func ShardTableNameByCrc32(tableName string, shardKey string, shardNum int) string {
	if shardNum == 0 {
		return ""
	}
	if tableName == "" || shardKey == "" {
		return ""
	}
	tableFormat := "_%04d" // 由于之前的表都是0000
	id, err := HashValue(shardKey)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(tableName+tableFormat, id%uint64(shardNum))
}

// 考虑直接使用hash := crc32.ChecksumIEEE([]byte(str)) >> 16 & 0x7FFF num := int(hash % 32)
func HashValue(value interface{}) (uint64, error) {
	switch val := value.(type) {
	case int:
		return uint64(val), nil
	case uint64:
		return val, nil
	case int64:
		return uint64(val), nil
	case string:
		if v, err := strconv.ParseUint(val, 10, 64); err != nil {
			return uint64(crc32.ChecksumIEEE(hack.Slice(val))), nil
		} else {
			return uint64(v), nil
		}
	case []byte:
		return uint64(crc32.ChecksumIEEE(val)), nil
	}
	return 0, errors.New("分表取hash失败类型不符合")
}
