package kid

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

// 时间戳纪元，从 2025-07-17T16:39:26.491+0800 开始计算
var Eopch int64 = 1752741566491

type ID [10]byte

const (
	// see https://github.com/kubernetes/kubernetes/blob/71b8ad965ef0e4cc2a2c22a68cf7db8d1f8c6dc8/staging/src/k8s.io/apimachinery/pkg/util/rand/rand.go#L83
	// 最大 24 个比特位
	machineAlphanums = "bcdfghjklmnpqrstvwxz2456789"
	timeMask         = 0x3FFFFFFFFFF
	encodingMask     = 0x1F
	machineMask      = 0xFFFFFF
	counterMask      = 0x3FFF
	lowbitMask       = 0x03
	nanoms           = 0xF4240
)

var (
	// 机器ID
	machineID int64
	// 自增ID
	counter uint32
	// k8s pod-name 名称正则
	hostnameReg = regexp.MustCompile(`-[a-z0-9]+-[` + machineAlphanums + `]{5}$`)
	// 32进制字典表
	encoding = "sh2tqf6lownky4pcg8xvze5jma9rd7bu"
)

var (
	ErrInvalidEncodingLength   = errors.New("invalid encoding length")
	ErrOnlyPrintableCharacters = errors.New("only printable characters")
	ErrRepeatedCharacters      = errors.New("repeated characters")
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("get hostname failed: %v", err))
	}
	machineID = getMachineID(hostname)
}

// getMachineID 没必要单独拆成方法，主要是为了测试
func getMachineID(hostname string) (id int64) {
	base := int64(len(machineAlphanums))
	charMap := map[rune]int{}
	for i, c := range machineAlphanums {
		charMap[c] = i
	}
	// 测试或正式环境 pod-name 肯定匹配
	if hostnameReg.MatchString(hostname) {
		for _, c := range hostname[len(hostname)-5:] {
			id = id*base + int64(charMap[c])
		}
		return
	}
	// 不匹配的属于本地测试，随机获取
	return rand.Int63n(machineMask)
}

// New 生成新的 kid
func New() ID {
	var id ID

	// 避免设置的 Eopch 太大或太小导致超过 42位
	now := ((time.Now().UnixNano() / nanoms) - Eopch) & 0x3FFFFFFFFFF
	// 时间戳占 42 位，大端序
	id[0] = byte(now >> 34)
	id[1] = byte(now >> 26)
	id[2] = byte(now >> 18)
	id[3] = byte(now >> 10)
	id[4] = byte(now >> 2)
	// 最后的2位放在下一个字节的7，8位上，所以这里要左移6位
	id[5] = byte((now & lowbitMask) << 6)

	// 机器ID, 24位
	// 这个字节的高2位已经被时间戳的最后2位占据，机器ID在这字节里只能存高6位
	// 24 - 6 = 18，向右移18位 再和高2位做 或 操作
	id[5] = id[5] | byte(machineID>>18)
	id[6] = byte(machineID >> 10)
	id[7] = byte(machineID >> 2)
	// 和时间戳的最后2位一样，这里也要左移6位
	id[8] = byte((machineID & lowbitMask) << 6)

	// 自增ID，14位(0~16383)
	i := atomic.AddUint32(&counter, 1) & counterMask
	// 这个字节高2位存放机器ID的最后2位，自增ID在这个字节只能存放高6位
	// 14 - 6 = 8，向右移8位 再和高2位做 或 操作
	id[8] = id[8] | byte(i>>8)
	id[9] = byte(i)

	return id
}

// SetMachineID 设置机器ID，0 ～ 16777215 之间
func SetMachineID(id uint64) {
	machineID = int64(id) & machineMask
}

// SetEncodingDict 设置编码字典，必须是 ascii 可见字符，长度32位
func SetEncodingDict(dict string) error {
	if len(dict) != 32 {
		return ErrInvalidEncodingLength
	}
	cs := make(map[rune]struct{}, 32)
	for _, c := range dict {
		if c < 33 || c > 126 {
			// 不可见的字符
			return ErrOnlyPrintableCharacters
		}
		if _, ok := cs[c]; ok {
			// 字典中出现重复字符
			return ErrRepeatedCharacters
		}
		cs[c] = struct{}{}
	}
	encoding = dict
	return nil
}

// String 获取自定义 base32 编码的ID 字符串
func (id ID) String() string {
	return string(id.encoding())
}

func (id ID) encoding() []byte {
	buf := make([]byte, 16)

	buf[0] = encoding[(id[0]>>3)&encodingMask]
	buf[1] = encoding[((id[0]<<2)|(id[1]>>6))&encodingMask]
	buf[2] = encoding[(id[1]>>1)&encodingMask]
	buf[3] = encoding[((id[1]<<4)|(id[2]>>4))&encodingMask]
	buf[4] = encoding[((id[2]<<4)|(id[3]>>7))&encodingMask]
	buf[5] = encoding[(id[3]>>2)&encodingMask]
	buf[6] = encoding[((id[3]<<3)|(id[4]>>5))&encodingMask]
	buf[7] = encoding[id[4]&encodingMask]
	buf[8] = encoding[(id[5]>>3)&encodingMask]
	buf[9] = encoding[((id[5]<<2)|(id[6]>>6))&encodingMask]
	buf[10] = encoding[(id[6]>>1)&encodingMask]
	buf[11] = encoding[((id[6]<<4)|(id[7]>>4))&encodingMask]
	buf[12] = encoding[((id[7]<<4)|(id[8]>>7))&encodingMask]
	buf[13] = encoding[(id[8]>>2)&encodingMask]
	buf[14] = encoding[((id[8]<<3)|(id[9]>>5))&encodingMask]
	buf[15] = encoding[id[9]&encodingMask]

	// 将前面空白部分的字符去掉(去掉前面的 00000)
	var start int
	for _, c := range buf {
		if c != encoding[0] {
			break
		}
		start++
	}
	if start > 0 {
		return buf[start:]
	}
	return buf
}

// Time 获取ID对应的毫秒级时间
func (id ID) Time() time.Time {
	t := (int64(id[0]) << 34) | (int64(id[1]) << 26) | (int64(id[2]) << 18) | (int64(id[3]) << 10) | (int64(id[4]) << 2) | (int64(id[5]) >> 6)

	// 如果 Eopch 在使用过程中被重新设置，获取的 Time 会和期望的不一致
	t += Eopch
	return time.Unix(t/1000, t%1000*nanoms)
}

// SequenceID 获取ID对应的自增序列号
func (id ID) SequenceID() int64 {
	return ((int64(id[8]) << 8) | int64(id[9])) & counterMask
}

// MachineID 获取机器ID
func (id ID) MachineID() int64 {
	return ((int64(id[5]) << 18) | (int64(id[6]) << 10) | (int64(id[7]) << 2) | int64(id[8])>>6) & machineMask
}

// MarshalText implements encoding/text TextMarshaler interface
func (id ID) MarshalText() ([]byte, error) {
	return id.encoding(), nil
}

// MarshalJSON returns a json byte array string of the snowflake ID.
func (id ID) MarshalJSON() ([]byte, error) {
	buf := id.encoding()
	result := make([]byte, len(buf)+2)
	result[0] = '"'
	copy(result[1:], buf)
	result[len(result)-1] = '"'

	return result, nil
}
