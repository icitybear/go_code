package loop_test

import (
	"sync"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

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

type rangeModel struct {
	Name string
	Age  int
}

func TestRange(t *testing.T) {
	// 循环里赋值
	list := []*rangeModel{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	for _, item := range list {
		item.Age += 1
	}
	// []rangeModel 如果非指针的情况下 [{a 1} {b 2} {c 3}]
	// []*rangeModel 指针情况下 [<*>{a 2} <*>{b 3} <*>{c 4}]

	// 协程竞争的情况
	wg := &sync.WaitGroup{}
	for _, item := range list {
		wg.Add(1)
		// 变量竞争的情况下 [<*>{a 2} <*>{b 3} <*>{gd_gd_c 4}]
		// go func() {
		// 	item.Name = spew.Sprintf("gd_%s", item.Name)
		// 	wg.Done()
		// }()

		// [<*>{gd_a 2} <*>{gd_b 3} <*>{gd_c 4}] 使用临时变量解决
		tmpItem := item
		go func() {
			tmpItem.Name = spew.Sprintf("gd_%s", tmpItem.Name)
			wg.Done()
		}()
	}
	wg.Wait()
	spew.Println(list)
}
