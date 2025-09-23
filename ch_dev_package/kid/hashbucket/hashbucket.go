package hashbucket

// GetBucketNum 计算字符串的哈希值并返回桶编号（范围 1 到 bucketAccount）
func GetBucketNum(str string, bucketAccount int) int64 {
	buf := []byte(str) // Go 字符串直接转换为 UTF-8 字节切片

	// 使用 uint64 初始化，避免溢出问题
	seed := uint64(0xcbf29ce484222325)
	for i := 0; i < len(buf); i++ {
		// 使用 uint64 进行位运算
		seed += (seed << 1) + (seed << 4) + (seed << 5) +
			(seed << 7) + (seed << 8) + (seed << 40)
		seed ^= uint64(buf[i])
	}

	// 将 uint64 转换为 int64 并取绝对值
	result := int64(seed)
	if result < 0 {
		result = -result
	}

	// 计算桶编号（范围 1 ~ bucketAccount）
	bucket := int64(bucketAccount)
	return result%bucket + 1
}

// GetBucketNumByConsistentHash 使用一致性哈希算法返回桶编号（范围 0 到 bucketAccount-1）
func GetBucketNumByConsistentHash(str string, bucketAccount int) int {
	buf := []byte(str)

	// 直接使用 uint64 类型
	seed := uint64(0xcbf29ce484222325)
	for i := 0; i < len(buf); i++ {
		seed += (seed << 1) + (seed << 4) + (seed << 5) +
			(seed << 7) + (seed << 8) + (seed << 40)
		seed ^= uint64(buf[i])
	}

	return consistentHash(seed, bucketAccount)
}

// 实现 Guava 的一致性哈希算法
func consistentHash(input uint64, buckets int) int {
	if buckets <= 0 {
		panic("buckets must be positive")
	}

	generator := newLinearCongruentialGenerator(input)
	candidate := 0

	for {
		// 计算下一个候选桶
		next := int(float64(candidate+1) / generator.nextDouble())
		if next >= 0 && next < buckets {
			candidate = next
		} else {
			return candidate
		}
	}
}

// 线性同余生成器
type linearCongruentialGenerator struct {
	state uint64
}

func newLinearCongruentialGenerator(seed uint64) *linearCongruentialGenerator {
	return &linearCongruentialGenerator{state: seed}
}

func (lcg *linearCongruentialGenerator) nextDouble() float64 {
	const multiplier = 2862933555777941757
	lcg.state = lcg.state*multiplier + 1

	// 无符号右移 33 位
	unsigned := lcg.state >> 33
	value := uint32(unsigned) // 取低 32 位

	// 转换为 [1, 0x80000000] 区间的浮点数
	return float64(value+1) / float64(0x80000000)
}
