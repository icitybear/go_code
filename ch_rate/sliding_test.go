package rate_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// https://github.com/CocaineCong/BiliBili-Code/blob/main/limit/sliding_counter_test.go
// SlidingWindowCounter 滑动窗口计数器限流器
type SlidingWindowCounter struct {
	limit     int64           // 限制数量
	window    time.Duration   // 时间窗口
	requests  map[int64]int64 // 时间戳到请求数的映射
	mutex     sync.Mutex      // 互斥锁
	precision time.Duration   // 精度（子窗口大小）
}

// NewSlidingWindowCounter 创建滑动窗口计数器
func NewSlidingWindowCounter(limit int64, window time.Duration, precision time.Duration) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		limit:     limit,
		window:    window,
		requests:  make(map[int64]int64),
		precision: precision,
	}
}

// Allow 检查是否允许请求通过
func (s *SlidingWindowCounter) Allow() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	now := time.Now()
	currentWindow := now.Truncate(s.precision).Unix()
	s.cleanExpiredWindows(now)                    // 清理过期的窗口数据
	totalRequests := s.countRequestsInWindow(now) // 计算当前窗口内的总请求数
	if totalRequests >= s.limit {                 // 检查是否超过限制
		return false
	}
	s.requests[currentWindow]++ // 增加当前窗口的计数
	return true
}

// cleanExpiredWindows 清理过期的窗口数据
func (s *SlidingWindowCounter) cleanExpiredWindows(now time.Time) {
	cutoff := now.Add(-s.window).Unix()
	for timestamp := range s.requests {
		if timestamp < cutoff {
			delete(s.requests, timestamp)
		}
	}
}

// countRequestsInWindow 计算窗口内的请求总数
func (s *SlidingWindowCounter) countRequestsInWindow(now time.Time) int64 {
	cutoff := now.Add(-s.window)
	total := int64(0)

	for timestamp, count := range s.requests {
		windowTime := time.Unix(timestamp, 0)
		if windowTime.After(cutoff) {
			total += count
		}
	}

	return total
}

// GetStatus 获取当前状态 返回窗口请求总数和当前限制数量
func (s *SlidingWindowCounter) GetStatus() (int64, int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	s.cleanExpiredWindows(now)
	current := s.countRequestsInWindow(now)

	return current, s.limit
}

// TestSlidingWindowCounter_Basic 测试滑动窗口计数器基本功能
func TestSlidingWindowCounter_Basic(t *testing.T) {
	limiter := NewSlidingWindowCounter(3, time.Second, 100*time.Millisecond)

	// 前3个请求应该通过
	for i := 0; i < 3; i++ {
		if limiter.Allow() {
			fmt.Printf("第%d个请求应该通过", i+1)
		}
	}

	// 第4个请求应该被拒绝
	if !limiter.Allow() {
		fmt.Printf("第4个请求应该通过")
	}

	// 检查状态
	current, limit := limiter.GetStatus()
	fmt.Printf("状态: current=%d, limit=%d", current, limit) // 3 3

}

// TestSlidingWindowCounter_SlidingWindow 测试滑动窗口特性
func TestSlidingWindowCounter_SlidingWindow(t *testing.T) {
	limiter := NewSlidingWindowCounter(2, 200*time.Millisecond, 50*time.Millisecond)

	// 发送2个请求
	limiter.Allow()
	limiter.Allow()

	// 应该被拒绝
	if !limiter.Allow() {
		fmt.Printf("第3个请求应该拒绝")
	}

	// 等待一半窗口时间，让部分请求过期
	time.Sleep(120 * time.Millisecond)

	// 现在应该可以通过一些请求
	if limiter.Allow() {
		fmt.Printf("滑动窗口后请求应该通过")
	}

	current, limit := limiter.GetStatus()
	fmt.Printf("状态: current=%d, limit=%d", current, limit) // 1 2
}

// TestSlidingWindowCounter_Precision 测试精度设置
func TestSlidingWindowCounter_Precision(t *testing.T) {
	// 高精度滑动窗口
	highPrecision := NewSlidingWindowCounter(5, time.Second, 10*time.Millisecond)

	// 低精度滑动窗口
	lowPrecision := NewSlidingWindowCounter(5, time.Second, 200*time.Millisecond)

	// 发送请求
	for i := 0; i < 3; i++ {
		highPrecision.Allow()
		lowPrecision.Allow()
	}

	// 检查状态
	highCurrent, _ := highPrecision.GetStatus()
	lowCurrent, _ := lowPrecision.GetStatus()

	fmt.Printf("精度测试: high=%d, low=%d", highCurrent, lowCurrent)

}

// TestSlidingWindowCounter_WindowExpiry 测试窗口过期
func TestSlidingWindowCounter_WindowExpiry(t *testing.T) {
	limiter := NewSlidingWindowCounter(3, 1000*time.Millisecond, 20*time.Millisecond)

	// 填满限制
	for i := 0; i < 3; i++ {
		limiter.Allow()
	}

	// 应该被拒绝
	if !limiter.Allow() {
		fmt.Printf("第4个请求应该被拒绝")
	}

	// 等待窗口完全过期
	time.Sleep(2000 * time.Millisecond)

	// 现在应该可以通过
	if limiter.Allow() {
		fmt.Printf("窗口过期后请求应该通过")
	}

	// 检查状态，应该只有1个请求 也就是上面的allow
	current, _ := limiter.GetStatus()
	fmt.Printf("窗口过期后计数为%d", current)

}

// TestSlidingWindowCounter_Concurrent 测试并发安全性
func TestSlidingWindowCounter_Concurrent(t *testing.T) {
	limiter := NewSlidingWindowCounter(100, time.Second, 10*time.Millisecond)
	var wg sync.WaitGroup
	var successCount int64
	var mu sync.Mutex

	// 启动200个并发请求
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// 应该只有100个请求成功
	if successCount != 100 {
		t.Errorf("期望100个成功请求，实际%d个", successCount)
	}
}

// 并发时同1秒超出了多个请求并非没有限制住
