package err_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
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

var BaseErr = errors.New("base error")

func TestErrWrap(t *testing.T) {
	err1 := fmt.Errorf("wrap base: %w", BaseErr)
	fmt.Println(err1)                          // wrap base: base error
	err2 := fmt.Errorf("wrap err1: %w", err1)  // %w 返回 wrapError类型
	fmt.Println(err2)                          // wrap err1: wrap base: base error
	err3 := fmt.Errorf("wrap err3: %+v", err1) // %w 返回 wrapError类型
	fmt.Println(err3)                          // wrap err3: wrap base: base error
}

func TestWrap(t *testing.T) {
	errA := errors.New("错误A")
	errB := fmt.Errorf("错误B: %w", errA)
	errC := fmt.Errorf("错误C: %w", errB)

	spew.Printf("errA: %+#v\nerrB: %+#v\nerrC: %+#v\n", errA, errB, errC)
}

func TestUnwrap(t *testing.T) {
	errA := errors.New("错误A")
	errB := fmt.Errorf("错误B: %w", errA)
	errC := fmt.Errorf("错误C: %w", errB)

	spew.Printf("errA: %+#v\nerrA Unwrap后: %+#v\n", errA, errors.Unwrap(errA))
	fmt.Println("---------------------------------------")
	spew.Printf("errB: %+#v\nerrB Unwrap后: %+#v\n", errB, errors.Unwrap(errB))
	fmt.Println("---------------------------------------")
	spew.Printf("errC: %+#v\nerrC Unwrap后: %+#v\nerrC两次Unwrap后: %+#v\n", errC, errors.Unwrap(errC), errors.Unwrap(errors.Unwrap(errC)))
}

func TestErrIs(t *testing.T) {
	err1 := fmt.Errorf("wrap base: %w", BaseErr)
	err2 := fmt.Errorf("wrap err1: %w", err1)
	println(err2 == BaseErr) // false
	if !errors.Is(err2, BaseErr) {
		panic("err2 is not BaseErr")
	}
	println("err2 is BaseErr") // err2 is BaseErr

	err := redis.Nil
	if errors.Is(err, redis.Nil) {
		println("err is redis.Nil")
	}
	// err3 := fmt.Errorf("wrap err3: %v", err) // err not is redis.Nil
	err3 := fmt.Errorf("wrap err3: %w", err) // err is redis.Nil
	if errors.Is(err3, redis.Nil) {
		println("err is redis.Nil")
	} else {
		println("err not is redis.Nil")
	}
}

// [errors.Is] 的核心逻辑：
// 功能： 判断一个错误 err 是否等于某个特定错误 target（或被其包裹）。
// 支持嵌套： 会递归地沿着 Unwrap() 或 Unwrap() []error 拆开 error 树，深度优先遍历。
// 匹配方式：
//
//	如果 err == target，返回 true；
//	如果 err 实现了 Is(error) bool 方法，并返回 true，也算匹配。
var ErrExample = errors.New("示例错误")

type ExampleErrorForIs struct{}

func (ce ExampleErrorForIs) Error() string {
	return "custom error"
}

func (ce ExampleErrorForIs) Is(target error) bool {
	return target.Error() == "custom error" // 自己实现 对比error的属性值
}
func TestErrorIs(t *testing.T) {
	errA := fmt.Errorf("包装第一层: %w", ErrExample)
	errB := fmt.Errorf("包装第二层: %w", errA)

	// 包装后的错误与原始错误直接比较：不相等
	// 包装后的错误，可通过errors.Is来判断是否是否个错误或由某个错误包装而来（深度优先遍历）
	if errors.Is(errA, ErrExample) {
		t.Log("errA 与 ErrExample 相等")
	} else {
		t.Log("errA 与 ErrExample 不相等")
	}
	if errors.Is(errB, errA) {
		t.Log("errB 与 errA 相等")
	} else {
		t.Log("errB 与 errA 不相等")
	}
	if errors.Is(errB, ErrExample) {
		t.Log("errB 与 ErrExample 相等")
	} else {
		t.Log("errB 与 ErrExample 不相等")
	}
}

// 实现了 Is 方法的自定义错误
func TestErrorIs_有实现Is方法(t *testing.T) {
	errA := errors.New("custom error")
	ce := &ExampleErrorForIs{}

	assert.True(t, errors.Is(ce, errA))
}

func TestErrorIs_相同错误内容的不同实例(t *testing.T) {
	var errMsg = "错误信息"
	errA := errors.New(errMsg)
	errB := errors.New(errMsg)

	spew.Printf("错误A: %+#v \n错误B: %+#v\n", errA, errB)

	// 直接比较：不相等（其实就是结构体比较）
	if errA == errB {
		fmt.Print("相等")
	}
	// Is判断是false
	if errors.Is(errA, errB) {
		fmt.Println("errA 与 errB 相等")
	} else {
		fmt.Println("errA 与 errB 不相等")
	}
	// 注意assert包的Equal是会判断为true的（因为会使用[reflect.DeepEqual]）
	assert.Equal(t, errA, errB) // 断言里会判断相等的

}

// 在错误链中查找第一个匹配某个类型的错误，如果找到则将其赋值给目标变量，并返回 true，否则返回 false。
// 🔍 匹配条件
// 错误的实际类型可以赋值给 target 指向的类型（即类型兼容）。
// 或者错误实现了 As(any) bool 方法，并返回 true。
//
// 🌲 错误树（Error Tree）
// 包括 err 本身；
// 以及通过不断调用 Unwrap() 或 Unwrap() []error 得到的嵌套错误；
// 多个嵌套错误时，采用**深度优先遍历（depth-first traversal）**方式查找。
//
// ⚠️ 注意事项
// target 必须是一个非 nil 的指针，指向实现了 error 接口的类型或接口类型，否则会 panic；
// 只能找到第一个匹配的错误。

type CustomError struct {
	Msg string
}

func (e *CustomError) Error() string {
	return e.Msg
}

// 值接收者实现 Error
type ValError struct {
	Msg string
}

func (e ValError) Error() string {
	return e.Msg
}

func TestErrorAs(t *testing.T) {
	ee := &CustomError{
		Msg: "a",
	}
	spew.Printf("ee: %+#v\n", ee) // ee: (*err_test.CustomError)(0x14000099190)a
	errA := fmt.Errorf("第一层包装: %w", ee)
	spew.Printf("errA: %+#v\n", errA) // errA: (*fmt.wrapError)(0x140000e22e0)第一层包装: a

	var ee1 *CustomError
	spew.Printf("ee1: %+#v\n&ee1: %+#v\n", ee1, &ee1)
	assert.Nil(t, ee1)                      // ee1: (*err_test.CustomError)<nil>
	assert.NotNil(t, &ee1)                  // &ee1: (**err_test.CustomError)(0x140000a0060)<nil>
	assert.True(t, errors.As(errA, &ee1))   // As执行 查找第一个匹配ee1类型的错误，如果找到（errA里找到匹配err1类型的）则将其赋值给目标变量err1，并返回 true
	spew.Printf("As操作后的ee1: %+#v\n", ee1)   // As操作后的ee1: (*err_test.CustomError)(0x14000099190)a
	spew.Printf("As操作后的&ee1: %+#v\n", &ee1) // As操作后的&ee1: (**err_test.CustomError)(0x140000a0060->0x14000099190)a
	fmt.Println("---------------------------")

	ee2 := &CustomError{}
	spew.Printf("ee2: %+#v\n&ee2: %+#v\n", ee2, &ee2)
	assert.NotNil(t, ee2)  // ee2: (*err_test.CustomError)(0x140000991e0) 空结构体 分配内存了 不为nil
	assert.NotNil(t, &ee2) // &ee2: (**err_test.CustomError)(0x140000a0068->0x140000991e0)
	assert.True(t, errors.As(errA, &ee2))
	spew.Printf("As操作后的ee2: %+#v\n", ee2)   // As操作后的ee2: (*err_test.CustomError)(0x14000099190)a
	spew.Printf("As操作后的&ee2: %+#v\n", &ee2) // As操作后的&ee2: (**err_test.CustomError)(0x140000a0068->0x14000099190)a
}

func TestAs_值接收者实现Error(t *testing.T) {
	// 使用 ValError（非指针）创建错误
	valErr := ValError{"this is a value error"}
	wrappedValErr := fmt.Errorf("wrap2: %w", valErr)

	spew.Printf("valErr: %+#v\nwrappedValErr: %+#v\n---------------\n", valErr, wrappedValErr)
	// valErr: (err_test.ValError)this is a value error
	// wrappedValErr: (*fmt.wrapError)(0x14000078300)wrap2: this is a value error

	var valTarget ValError
	spew.Printf("valTarget: %+#v\n&valTarget: %+#v\n", valTarget, &valTarget)
	// valTarget: (err_test.ValError)
	// &valTarget: (*err_test.ValError)(0x140000671f0)
	assert.True(t, errors.As(wrappedValErr, &valTarget))         // 此时找到
	spew.Printf("转换后的valTarget: %+#v\n-----------\n", valTarget) // 转换后的valTarget: (err_test.ValError)this is a value error

	var valTargetPtr *ValError // 使用指针
	spew.Printf("valTargetPtr: %+#v\n&valTargetPtr: %+#v\n", valTargetPtr, &valTargetPtr)
	// valTargetPtr: (*err_test.ValError)<nil>
	// &valTargetPtr: (**err_test.ValError)(0x14000010070)<nil>
	assert.False(t, errors.As(wrappedValErr, &valTargetPtr)) // 此时没找到

	spew.Printf("转换后的valTargetPtr: %+#v\n", valTargetPtr)
	// 转换后的valTargetPtr: (*err_test.ValError)<nil>
}

type CustomError2 struct {
	Code int
	Msg  string
}

func (c *CustomError2) Error() string {
	return fmt.Sprintf("自定义错误, code: %d, msg: %s", c.Code, c.Msg)
}
func Handle(content string) error {
	var e *CustomError2
	if content != "" {
		return e // 容易踩坑的点
	}
	e = &CustomError2{
		Code: 1,
		Msg:  "内容为空",
	}
	return e
}
func TestError(t *testing.T) {
	c1 := "xx"
	err := Handle(c1)
	// 踩坑点 容易直接判断为nil 此时 e 是一个 error接口变量，它的动态类型是 *CustomError2，动态值是 nil。
	// 由于动态类型不为 nil，因此 e == nil 的结果是 false。
	if err != nil {
		spew.Printf("处理c1失败, err=%+#v\n", err) // 处理c1失败, err=(*err_test.CustomError2)<nil>
	}
	c2 := ""
	err = Handle(c2)
	if err != nil {
		spew.Printf("处理c2失败, err=%+#v\n", err) // 处理c2失败, err=(*err_test.CustomError2)(0x1400000e300)自定义错误, code: 1, msg: 内容为空
	}
}

// 1.20版本才有方法errors.Join
// func TestJoin(t *testing.T) {
// 	errA := errors.New("错误A")
// 	errB := errors.New("错误B")
// 	errC := &CustomError{Msg: "错误C"}

// 	errJoined := errors.Join(errA, errB, errC)

// 	spew.Printf("errJoined: %+#v\n", errJoined)
// 	assert.True(t, errors.Is(errJoined, errA))
// 	assert.True(t, errors.Is(errJoined, errB))
// 	assert.True(t, errors.Is(errJoined, errC))
// 	var ce *CustomError
// 	assert.True(t, errors.As(errJoined, &ce))
// }
