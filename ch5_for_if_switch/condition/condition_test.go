package condition_test

import (
	"testing"
)

// 多个条件 正确的写法
func TestSwitchMultiCase(t *testing.T) {
	for i := 0; i < 5; i++ {
		switch i {
		case 0, 2:
			t.Log("Even")
		case 1, 3:
			t.Log("Odd")
		default:
			t.Log("it is not 0-3")
		}
	}
}

func TestSwitchCaseCondition2(t *testing.T) {
	i := 3
	switch i {
	case 3: // 并不会执行到下面的case
	case 1:
		t.Log(i)
	default:
		t.Log("unknow")
	}
	t.Log("after switch")
}

// 表达式
func TestSwitchCaseCondition(t *testing.T) {
	for i := 0; i < 5; i++ {
		switch {
		case i%2 == 0:
			t.Log("Even")
		case i%2 == 1:
			t.Log("Odd")
		default:
			t.Log("unknow")
		}
	}
}
