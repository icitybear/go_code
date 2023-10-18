package once_test

import (
	"fmt"
	"sync"
	"testing"
	"unsafe"
)

type Singleton struct {
	data string
}

var singleInstance *Singleton
var once sync.Once

func GetSingletonObj() *Singleton {
	//仅执行一次
	once.Do(func() {
		fmt.Println("Create Obj")
		singleInstance = new(Singleton)
	})
	return singleInstance
}

func TestGetSingletonObj(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		//调用了10次
		wg.Add(1)
		go func() {
			obj := GetSingletonObj()                //只有1次Create Obj
			fmt.Printf("%X\n", unsafe.Pointer(obj)) //指针地址都是同一个 因为是单例模式
			wg.Done()
		}()
	}
	wg.Wait()
}
