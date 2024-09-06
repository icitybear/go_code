package context_cancel

import (
	"fmt"
	"testing"
)

type ctx struct {
	name string
	idx  int8
	h    handlers
}

/*
*
idx=0 func1 start

	idx = 1 func2 start
		idx = 2 func3 start
		idx = 3 return
	idx =1 func2 end

idx=0 func1 end
*
*/
func (c *ctx) next() {
	c.idx++
	fmt.Printf("idx:%d\n", c.idx)
	if c.idx < int8(len(c.h)) {
		c.h[c.idx](c) // 核心调用
		c.idx++
	}
}

type handler func(*ctx) // 适配器 没有返回值的函数类型, 也没有函数体

type handlers []handler // 函数数组

type engin struct {
	name string
	hds  handlers
}

func (e *engin) run(c *ctx) {
	fmt.Printf("engine:%s running\n", e.name)
	c.h = e.hds
	c.next() // ctx往注册的继续调用
	fmt.Printf("engine:%s end\n", e.name)
}

// 原理-大概理解
func TestYl(t *testing.T) {
	testHandler()
}

func testHandler() {

	h1 := func(c *ctx) {
		fmt.Println("h1 start")
		c.next()
		fmt.Println("h1 end")
	}

	h2 := func(c *ctx) {
		fmt.Println("h2 start")
		c.next()
		fmt.Println("h2 end")
	}

	h3 := func(c *ctx) {
		fmt.Println("h3 running...") // 没继续执行 c.next()
	}

	hds := make([]handler, 3)
	// 按顺序排好处理函数
	hds[0] = h1
	hds[1] = h2
	hds[2] = h3

	c1 := &ctx{
		idx:  -1,
		name: "ctx1",
	}

	engin1 := &engin{
		name: "engin1",
		hds:  hds,
	}
	engin1.run(c1)
	// engine:engin1 running
	// idx:0
	// h1 start
	//
	//	idx:1
	//	h2 start
	//		idx:2
	//		h3 running...
	//	h2 end
	//
	// h1 end
	// engine:engin1 end
}
