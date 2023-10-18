package constant_test

import "testing"

const Pi float64 = 3.14159265358979323846
const zero = 0.0 // 无类型浮点常量
const (          // 通过一个 const 关键字定义多个常量，和 var 类似
	size int64 = 1024
	eof        = -1 // 无类型整型常量
)
const u, v float32 = 0, 3   // u = 0.0, v = 3.0，常量的多重赋值
const a, b, c = 3, 4, "foo" // a = 3, b = 4, c = "foo", 无类型整型和字符串常量

const (
	Monday = 1 + iota //+1 iota一开始默认0
	Tuesday
	Wednesday
)

// 在每一个 const 关键字出现时被重置为 0，然后在下一个 const 出现之前，每出现一次 iota，其所代表的数字会自动增 1
// 省略后一个赋值表达式(相同的表达式)
const (
	Readable = 1 << iota //左移位
	Writable
	Executable
)

func TestConstantTry(t *testing.T) {
	t.Log(Monday, Tuesday)
}

func TestConstantTry1(t *testing.T) {
	t.Log(Readable, Writable, Executable)
	a := 1 //0001
	t.Log(a&Readable == Readable, a&Writable == Writable, a&Executable == Executable)
}

const (
	mutexLocked      = 1 << iota // 1左移0位（1） mutex is locked 上锁 1 state右侧的第一个 bit位标志是否上锁，0-未上锁，1-已上锁；
	mutexWoken                   // 1左移1位 2 state右侧的第二个bit位标 是否有唤醒的在抢占 是否有 goroutine 从阻塞中被唤醒，0-没有，1-有；
	mutexStarving                // 1左移2位 4 state右侧的第三个bit位标 饥渴状态
	mutexWaiterShift = iota      // iota=3 用来移位3位

	starvationThresholdNs = 1e6 // 1nm，自旋超过时间阀值就从正常模式切换到饥渴模式
)

func TestConstant(t *testing.T) {
	t.Log(mutexLocked, mutexWoken, mutexStarving, mutexWaiterShift)
}
