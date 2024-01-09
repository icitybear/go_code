package new_package

import (
	"fmt"
	"series" // 使用本地开发的包sdk 一个模块（module new_package）包含多个包
	"testing"
)

// 这里改成main包 那就是入口了 go.mod sum讲解看刘丹冰
// mod可能模块或者包间接引入
func TestX(t *testing.T) {
	fmt.Println("Hello World")
	num := series.Square(5)
	fmt.Println(num)
}

// module方式解决 GOPATH的弊端 1⽆版本控制概念 2⽆法同步⼀致第三⽅版本号 3
