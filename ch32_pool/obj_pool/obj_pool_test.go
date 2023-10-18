package object_pool

//用的同一个包名 所以方法变量共用

import (
	"fmt"
	"testing"
	"time"
)

func TestObjPool(t *testing.T) {
	pool := NewObjPool(10)
	// 如果先池子生成10个对象，再额外放新的 会报错
	// if err := pool.ReleaseObj(&ReusableObj{}); err != nil { //尝试放置超出池大小的对象
	// 	t.Error(err)
	// }

	for i := 0; i < 11; i++ {
		if v, err := pool.GetObj(time.Second * 1); err != nil {
			t.Error(err)
		} else {
			fmt.Printf("%T\n", v)
			if err := pool.ReleaseObj(v); err != nil {
				t.Error(err)
			}
		}

	}

	fmt.Println("Done")
}
