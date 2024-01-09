package fugaitest

import (
	"fmt"
	"testing"
)

// 变量的生命周期 作用域
func TestSM(t *testing.T) {
	var a int
	if c := 3; 1 > 0 {
		// a = 1 // 这样就是直接使用外部的了

		a := 2         // 这a重新定义了
		fmt.Println(a) // 输出2 a是这if作用域里的 而不是外部的
	} else {
		fmt.Println(c) //
	}
	fmt.Println(a) // 输出0 因为没经过
}
