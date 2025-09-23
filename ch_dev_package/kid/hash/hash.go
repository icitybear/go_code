package hash

import (
	"encoding/binary"
)

// github.com/spaolacci/murmur3 64为版本直接用这个
//
//	Go 语言中实现 MurmurHash3 算法（32-bit 版本） MurmurHash算法的不同变体（例如32位、64位、128位版本）需要不同的实现。
func Murmur3_32(key []byte, seed uint32) uint32 {
	const (
		c1 uint32 = 0xcc9e2d51
		c2 uint32 = 0x1b873593
		r1 uint32 = 15
		r2 uint32 = 13
		m  uint32 = 5
		n  uint32 = 0xe6546b64
	)

	hash := seed
	length := len(key)
	roundedEnd := (length & 0xfffffffc) // round down to 4 byte block

	for i := 0; i < roundedEnd; i += 4 {
		k1 := binary.LittleEndian.Uint32(key[i : i+4])
		k1 *= c1
		k1 = (k1 << r1) | (k1 >> (32 - r1))
		k1 *= c2

		hash ^= k1
		hash = ((hash<<r2)|(hash>>(32-r2)))*m + n
	}

	if length > roundedEnd {
		tail := uint32(0)
		switch length & 3 {
		case 3:
			tail |= uint32(key[roundedEnd+2]) << 16
			fallthrough
		case 2:
			tail |= uint32(key[roundedEnd+1]) << 8
			fallthrough
		case 1:
			tail |= uint32(key[roundedEnd])
			tail *= c1
			tail = (tail << r1) | (tail >> (32 - r1))
			tail *= c2
			hash ^= tail
		}
	}

	hash ^= uint32(length)
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16

	return hash
}
