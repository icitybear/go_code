package rate_test

import (
	"sync"
	"time"
)

type TimeSlice struct {
	timestamp int64 // 精确到毫秒的时间标识
	count     int64 // 该毫秒内的请求计数
}

// 滑动窗口限流器
type SlidingWindowLimiter1 struct {
	windowDuration int64        // 窗口总时长(单位：毫秒)
	maxRequests    int64        // 窗口内允许的最大请求数
	slices         []*TimeSlice // 动态时间片集合（非严格环形）
	currentPos     int          // 最新时间片写入位置指针
	mutex          sync.Mutex   // 并发控制锁
}

// SlidingWindowLimiter 创建滑动窗口计数器
func NewSlidingWindowLimiter1(limit int64, window int64) *SlidingWindowLimiter1 {
	return &SlidingWindowLimiter1{
		maxRequests:    limit,
		windowDuration: window,
		mutex:          sync.Mutex{},
		slices:         make([]*TimeSlice, 0, window/limit),
	}
}

func (l *SlidingWindowLimiter1) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now().UnixNano() / 1e6 // 当前毫秒时间戳
	// 清理过期时间片
	cutoff := now - l.windowDuration // 过期的毫秒时间戳
	for i := 0; i < len(l.slices); {
		if l.slices[i].timestamp < cutoff {
			l.slices = append(l.slices[:i], l.slices[i+1:]...)
		} else {
			i++
		}
	}

	// 统计当前窗口请求总数
	var total int64
	for _, s := range l.slices {
		total += s.count
	}
	if total >= l.maxRequests {
		return false // 大于最多请求数
	}

	// tag: 在清理过期时间片后，l.currentPos可能指向无效位置
	// 当l.currentPos位置被清理后，访问l.slices[l.currentPos]会导致索引越界

	// 记录当前请求
	if len(l.slices) > 0 && l.slices[l.currentPos].timestamp == now {
		l.slices[l.currentPos].count++
	} else {
		l.slices = append(l.slices, &TimeSlice{
			timestamp: now,
			count:     1,
		})
		l.currentPos = len(l.slices) - 1
	}
	return true
}

// 优化版
// tag: 解决索引越界问题：
// 移除了currentPos指针，改用时间片列表尾部作为最新时间片
// 使用head和tail指针管理有效时间片范围
// tag: 性能优化：
// 添加totalRequests字段缓存总请求数，避免每次遍历计算
// 使用头指针head标记有效数据起始点，避免频繁删除操作
// 当头部空闲空间过大时，压缩切片减少内存占用
// 内存优化：
// 初始化时预分配切片容量，减少动态扩容
// 使用指针切片([]*TimeSlice)减少内存复制
type SlidingWindowLimiter struct {
	windowDuration int64        // 窗口总时长(单位：毫秒)
	maxRequests    int64        // 窗口内允许的最大请求数
	slices         []*TimeSlice // 时间片集合（按时间排序）
	head           int          // 有效时间片的起始索引
	// tail           int          // 有效时间片的结束索引
	totalRequests int64      // 当前窗口内总请求数
	mutex         sync.Mutex // 并发控制锁
}

func NewSlidingWindowLimiter(limit int64, window int64, precision int64) *SlidingWindowLimiter {
	if precision <= 0 {
		precision = 1
	}

	return &SlidingWindowLimiter{
		maxRequests:    limit,
		windowDuration: window,
		slices:         make([]*TimeSlice, 0, window/precision), // 预分配空间
		mutex:          sync.Mutex{},
	}
}

func (l *SlidingWindowLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now().UnixNano() / 1e6 // 当前毫秒时间戳
	cutoff := now - l.windowDuration   // 过期时间点

	// 清理过期时间片（从头部开始）
	for l.head < len(l.slices) {
		if l.slices[l.head].timestamp >= cutoff {
			break
		}
		l.totalRequests -= l.slices[l.head].count
		l.head++
	}

	// 如果头部指针超过切片一半，则压缩切片
	if l.head > len(l.slices)/2 {
		l.slices = append(make([]*TimeSlice, 0, cap(l.slices)), l.slices[l.head:]...)
		l.head = 0
	}

	// 检查是否超过限制
	if l.totalRequests >= l.maxRequests {
		return false
	}

	// 尝试更新最近的时间片
	if len(l.slices) > l.head {
		last := l.slices[len(l.slices)-1]
		if last.timestamp == now {
			last.count++
			l.totalRequests++
			return true
		}
	}

	// 添加新的时间片
	l.slices = append(l.slices, &TimeSlice{
		timestamp: now,
		count:     1,
	})
	l.totalRequests++
	return true
}
