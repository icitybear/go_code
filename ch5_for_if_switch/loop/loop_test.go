package loop_test

import "testing"

func TestWhileLoop(t *testing.T) {
	n := 0
	for n < 5 {
		t.Log(n)
		n++
	}
}

func TestForSelect(t *testing.T) {
	// for配合select break只是跳出select, 继续下次循环， 相当continue
	// for配合switch break只是跳出switch, 继续下次循环
	for i := 0; i < 5; i++ {
		switch {
		case i%2 == 0:
			t.Log("Even")
		case i%2 == 1:
			t.Log("Odd")
			break
			t.Log("hhh")
		default:
			t.Log("unknow")
		}
	}
}
