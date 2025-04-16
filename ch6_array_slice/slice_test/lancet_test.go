package slice_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
)

// tag: 1.21版本编译通不过，1.22可以 建议使用1.22
// type declarations inside generic functions are not currently supported 这个包2.3.3版本有问题，直接用2.3.5
// func Join[T any](s []T, separator string) string 范型参数的隐式转换

func TestXxx2(t *testing.T) {
	// arr := []int32{1, 2, 3}
	arr := []int32{}
	str := slice.Join(arr, ",")
	fmt.Println(str, len(str)) // 0
	// 1,2,3  5个字符串

	// res0 := strutil.SplitAndTrim("", ",") // 避免了直接使用  strings.Split问题  每个元素都trim
	// fmt.Println(res0, len(res0))            // [] 0
	res0 := strutil.SplitEx("", ",", false) // 是否去除空字符串
	fmt.Println(res0, len(res0))            // [] 1
	res1 := strutil.SplitEx("", ",", true)  // 是否去除空字符串
	fmt.Println(res1, len(res1))            // [] 0
}

func TestXxx11(t *testing.T) {
	// str := Int32Join([]int32{1, 2, 3}, ",")
	str := Int32Join([]int32{}, ",")
	fmt.Println(str, len(str)) // 0
	// 1,2,3  5个字符串

	fmt.Println(StrArrToInt32Arr(str, ","))

	// tag: 区别
	res0 := strings.Split("", ",")
	fmt.Println(res0, len(res0))     // [] 1
	res := StrArrToInt32Arr("", ",") // 如果直接使用 strings.Split
	fmt.Println(res, len(res))       // [] 0
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

func TestSplitTrim(t *testing.T) {
	str := "{a,b,  c, d}"
	res := strutil.SplitAndTrim(str, ",", "{}")
	fmt.Println(res)
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
