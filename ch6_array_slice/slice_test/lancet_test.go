package slice_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// tag: 1.21版本编译通不过，1.22可以 建议使用1.22
// func Join[T any](s []T, separator string) string 范型参数的隐式转换

// func TestXxx9(t *testing.T) {
// 	nums := []int{1, 2, 3, 4, 5}
// 	result1 := slice.Join(nums, ",") // 范型
// 	fmt.Println(result1)

// }

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

// func TestXxx11(t *testing.T) {
// 	result1 := slice.AppendIfAbsent([]string{"a", "b"}, "b")
// 	result2 := slice.AppendIfAbsent([]string{"a", "b"}, "c")

// 	fmt.Println(result1)
// 	fmt.Println(result2)
// }

func TestXxx11(t *testing.T) {
	// str := Int32Join([]int32{1, 2, 3}, ",")
	str := Int32Join([]int32{}, ",")
	fmt.Println(str, len(str))

	fmt.Println(StrArrToInt32Arr(str, ","))
}
func Int32Join(arr []int32, str string) string {
	var strArr []string
	for _, v := range arr {
		strArr = append(strArr, fmt.Sprintf("%d", v))
	}
	return strings.Join(strArr, str)
}

func StrArrToInt32Arr(str string, sep string) []int32 {

	var intList []int32
	if len(str) == 0 {
		return intList
	}
	arr := strings.Split(str, sep) // 缺点: 至少会有一个 就算是空字符串

	for _, v := range arr {
		val, _ := strconv.Atoi(v)
		intList = append(intList, int32(val))
	}
	return intList
}

func TestChunks(t *testing.T) {
	strArr := []string{"aa", "bb", "cc"}
	fmt.Println(strArr)
	result := SplitIntoChunks(strArr, 2)

	fmt.Println(result[0])
	fmt.Println(result[1])
}

func SplitIntoChunks[T any](s []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
