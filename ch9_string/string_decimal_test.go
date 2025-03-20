package string_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

// 字符串与数值转字符串
func TestAoti(t *testing.T) {
	str := "7431820065157906458"
	// string转int     Itoa
	s, _ := strconv.Atoi(str) // 整形最大范围 9223372036854775807
	fmt.Println(s)            // 7431820065157906458

	// 指数形式表示的浮点型float64
	floatNum := 7.431820065157906e+18
	// float64转string
	strNum := strconv.FormatFloat(floatNum, 'f', 15, 64)  // 'f' 表示没有指数部分，保留15位小数 一般会是2保留2位
	fmt.Println(strNum)                                   // 7431820065157906432.000000000000000
	strNum1 := strconv.FormatFloat(floatNum, 'f', -1, 64) // 'f' 表示去掉指数部分 -1没有小数位
	fmt.Println(strNum1)                                  // 7431820065157906000
	// int64转string
	// strconv.FormatInt(int64(num1), 10) // 10进制

	// string转float64 float32
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", floatNum), 64) // 参数只有32或者64
	fmt.Println(num)                                                // 输出原始的 float64 数值 7.431820065157906e+18

	// 使用 decimal.NewFromFloat 创建一个 decimal.Decimal 实例
	decimalValue := decimal.NewFromFloat(floatNum)
	// 乘以 100，使用 decimal 的 Mul 方法
	decimalValue = decimalValue.Mul(decimal.NewFromInt(100))

	// 将结果转换回 float64
	res, _ := decimalValue.Float64()
	fmt.Println(res) // 7.431820065157906e+20 因为多了2位

	// string转decimal 7431820065157906000"
	d, _ := decimal.NewFromString(strNum1)
	fmt.Println(d.String()) // 转string "7431820065157906000"
	res2, _ := d.Float64()  // 转float64 7431820065157906000
	fmt.Println(strconv.FormatFloat(res2, 'f', -1, 64))

}

// 通用转换方式  int也适用
func GetInterFaceDecimal(v interface{}) decimal.Decimal {
	var r decimal.Decimal
	var err error
	switch v.(type) {
	case uint:
		r = decimal.NewFromInt(int64(v.(uint)))
		break
	case int8:
		r = decimal.NewFromInt(int64(v.(int8)))
		break
	case uint8:
		r = decimal.NewFromInt(int64(v.(uint8)))
		break
	case int16:
		r = decimal.NewFromInt(int64(v.(int16)))
		break
	case uint16:
		r = decimal.NewFromInt(int64(v.(uint16)))
		break
	case int32:
		r = decimal.NewFromInt32(v.(int32))
		break
	case uint32:
		r = decimal.NewFromInt(int64(v.(uint32)))
		break
	case int64:
		r = decimal.NewFromInt(v.(int64))
		break
	case uint64:
		r = decimal.NewFromInt(int64(v.(uint64)))
		break
	case float32:
		r = decimal.NewFromFloat32(v.(float32))
		break
	case float64:
		r = decimal.NewFromFloat(v.(float64))
		break
	case string:
		r, err = decimal.NewFromString(v.(string))
		if err != nil {
			return decimal.NewFromInt(0)
		}
		break
	case int:
		r = decimal.NewFromInt(int64(v.(int)))
		break
	case nil:
		return decimal.NewFromInt(0)
	case decimal.Decimal:
		r = v.(decimal.Decimal)
	default:
		return decimal.NewFromInt(0)
	}
	return r
}

func DivideFromInterface(x interface{}, y interface{}) decimal.Decimal {
	xd := GetInterFaceDecimal(x)
	yd := GetInterFaceDecimal(y)
	rs := xd.Div(yd)
	return rs
}

func MulFromInterface(x interface{}, y interface{}) decimal.Decimal {
	xd := GetInterFaceDecimal(x)
	yd := GetInterFaceDecimal(y)
	rs := xd.Mul(yd)
	return rs
}

// string转千分计数且考虑小数点
func ThousandSeparate(str string) string {
	arr := strings.Split(str, ".")
	integerPart := arr[0]
	var result strings.Builder

	// 处理整数部分
	for i, ch := range integerPart {
		// 从右往左每三位加逗号（当剩余位数是3的倍数且不在起始位置时）
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			result.WriteByte(',')
		}
		result.WriteRune(ch)
	}

	// 处理小数部分
	if len(arr) > 1 {
		result.WriteByte('.')
		result.WriteString(strings.Join(arr[1:], "."))
	}

	return result.String()
}
