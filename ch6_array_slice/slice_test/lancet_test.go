package slice_test

import (
	"fmt"

	"testing"

	// "github.com/duke-git/lancet"
	"github.com/duke-git/lancet/v2/slice"
)

// tag: 1.21版本编译通不过，1.22可以 建议使用1.22
// func Join[T any](s []T, separator string) string 范型参数的隐式转换

func TestXxx9(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	result1 := slice.Join(nums, ",") // 范型
	fmt.Println(result1)

}

// func TestXxx10(t *testing.T) {
// 	str := ",1,2,4,"
// 	permission := strings.Split(str, ",") //
// 	for k, v := range permission {
// 		if v == "" {
// 			continue
// 		}
// 		fmt.Printf("k=%v v=%v\n", k, v)
// 	}
//
// 	// 范型
// 	filtered := lancet.Filter(permission, func(v int) bool {
// 		return v != 0
// 	})
// 	for k, v := range filtered {
// 		fmt.Printf("k=%v v=%v\n", k, v)
// 	}
// }

func TestXxx11(t *testing.T) {
	result1 := slice.AppendIfAbsent([]string{"a", "b"}, "b")
	result2 := slice.AppendIfAbsent([]string{"a", "b"}, "c")

	fmt.Println(result1)
	fmt.Println(result2)
}
