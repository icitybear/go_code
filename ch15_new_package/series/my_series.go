package series

import "fmt"

// import > const > var > init
// init 初始化 上到下 内到外
func init() {
	fmt.Println("init1")
}

func init() {
	fmt.Println("init2")
}

func Square(n int) int {
	return n * n
}

func GetFibonacciSerie(n int) []int {
	ret := []int{1, 1}
	for i := 2; i < n; i++ {
		ret = append(ret, ret[i-2]+ret[i-1])
	}
	return ret
}
