package object_pool

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

//协程安全 有锁的开销  锁的代价大还是对象创建的代价大，GC  复用
//sync.Pool 对象池  对象池的作用就是缓存对象，减少对象创建的开销，提高性能。 获取与放回
func TestSyncPool(t *testing.T) {
	pool := &sync.Pool{
		//指定New
		New: func() interface{} {
			fmt.Println("Create a new object.")
			return 100
		},
	}

	v := pool.Get().(int) //一开始没私有对象 共享也没有 去New了 并返回100
	fmt.Println(v)
	fmt.Println("put before")
	//pool.Put(33)
	runtime.GC()              //GC 会清除sync.pool中缓存的对象
	v1, _ := pool.Get().(int) //只要不put 放回 那么私有对象就没放 再放就是到P下的共享池
	fmt.Println(v1)
	//pool.Put(55) //因为有私有对象了  就不会再去new， 直接放心值55
	v2, _ := pool.Get().(int) //Get 私有对象 P下的共享池 其他P的共享池
	fmt.Println(v2)
}

func TestSyncPoolInMultiGroutine(t *testing.T) {
	//多协程下  对象缓存的 工作机制  不同Processor私有对象 共享对象
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Create a new object2.")
			return 10
		},
	}
	//先放入3个100
	pool.Put(100)
	pool.Put(100)
	pool.Put(100)

	//多个协程情况下 协程安全 私有对象 共享池 时有锁的开销
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			//只要不put 放回 那么私有对象就没放 Get都是会New
			fmt.Println(pool.Get()) //3个100 7个New
			wg.Done()
		}(i)
	}
	wg.Wait()
}
