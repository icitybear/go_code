package err_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

var LessThanTwoError = errors.New("n should be not less than 2")
var LargerThenHundredError = errors.New("n should be not larger than 100")

// fmt.Errorf("用户:%v或设备:%s未找到归因数据", tokenUuid[token], token) 替代 errors.New()

// 实现error接口就行
// 现在基本用 errors.Is 和errors.As 2方法来判断错误具体归属
//
// 函数有返回错误error接口 有错误要处理
func GetFibonacci(n int) ([]int, error) {
	if n < 2 {
		// 切片 nil 和具体的错误接口实例指针
		return nil, LessThanTwoError
	}
	if n > 100 {
		return nil, LargerThenHundredError
	}
	fibList := []int{1, 1}

	for i := 2; /*短变量声明 := */ i < n; i++ {
		fibList = append(fibList, fibList[i-2]+fibList[i-1])
	}
	return fibList, nil
}

// 错误处理
func TestGetFibonacci(t *testing.T) {
	if v, err := GetFibonacci(1); err != nil {
		if err == LessThanTwoError {
			fmt.Println("It is less. LessThanTwoError")
		}
		t.Log(err) //t.Error(err) 直接报错
	} else {
		t.Log(v)
	}

}

func GetFibonacci1(str string) {
	var (
		i    int
		err  error
		list []int
	)
	if i, err = strconv.Atoi(str); err == nil {
		if list, err = GetFibonacci(i); err == nil {
			fmt.Println(list)
		} else {
			fmt.Println("Error", err)
		}
	} else {
		fmt.Println("Error", err)
	}
}

func GetFibonacci2(str string) {
	var (
		i    int
		err  error
		list []int
	)
	//正常写法 应该是有异常错误 就尽早处理终止
	if i, err = strconv.Atoi(str); err != nil {
		fmt.Println("Error", err)
		return
	}
	if list, err = GetFibonacci(i); err != nil {

		fmt.Println("Error", err)
		return
	}
	fmt.Println(list)

}

var BaseErr = errors.New("the underlying base error")

func TestErrWrap(t *testing.T) {
	err1 := fmt.Errorf("wrap base: %w", BaseErr)
	fmt.Println(err1)
	err2 := fmt.Errorf("wrap err1: %w", err1) // %w 返回 wrapError类型
	fmt.Println(err2)
}

// 打印结果：
// wrap base: the underlying base error
// wrap err1: wrap base: the underlying base error

func TestErrIs(t *testing.T) {
	err1 := fmt.Errorf("wrap base: %w", BaseErr)
	err2 := fmt.Errorf("wrap err1: %w", err1)
	println(err2 == BaseErr) // false
	if !errors.Is(err2, BaseErr) {
		panic("err2 is not BaseErr")
	}
	println("err2 is BaseErr")
}

// 打印结果：
// false
// err2 is BaseErr

type TypicalErr struct {
	e string
}

func (t TypicalErr) Error() string {
	return t.e
}

func TestErrAs(t *testing.T) {
	err := TypicalErr{"typical error"}
	err1 := fmt.Errorf("wrap err: %w", err)
	err2 := fmt.Errorf("wrap err1: %w", err1)
	var e TypicalErr
	if !errors.As(err2, &e) {
		panic("TypicalErr is not on the chain of err2")
	}
	println("TypicalErr is on the chain of err2")
	println(err == e)
}

// 打印结果：
// TypicalErr is on the chain of err2
// true
