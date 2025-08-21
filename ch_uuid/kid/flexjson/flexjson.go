package flexjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/bytedance/sonic"
)

// FlexInt 将 int 或 字符串数字解析为 int
type FlexInt int

// FlexInt8 将 int8 或 字符串数字解析为 int8
type FlexInt8 int8

// FlexInt16 将 int16 或 字符串数字解析为 int16
type FlexInt16 int16

// FlexInt32 将 int32 或 字符串数字解析为 int32
type FlexInt32 int32

// FlexInt64 将 int64 或 字符串数字解析为 int64
type FlexInt64 int64

// FlexUint 将 uint 或 字符串数字解析为 uint
type FlexUint uint

// FlexUint8 将 uint8 或 字符串数字解析为 uint8
type FlexUint8 uint8

// FlexUint16 将 uint16 或 字符串数字解析为 uint16
type FlexUint16 uint16

// FlexUint32 将 uint32 或 字符串数字解析为 uint32
type FlexUint32 uint32

// FlexUint64 将 uint64 或 字符串数字解析为 uint64
type FlexUint64 uint64

// FlexFloat32 将 float32 或 字符串数字解析为 float32
type FlexFloat32 float32

// FlexFloat64 将 float64 或 字符串数字解析为 float64
type FlexFloat64 float64

// FlexBool 将 1,'1','T','t','true','TRUE','True' 解析为 true，其他解析为 false
type FlexBool bool

// FlexObject 兼容 PHP 对象为空时可能返回 [] 的情况
type FlexObject[T any] struct {
	ptr *T
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexInt) UnmarshalJSON(b []byte) error {
	value, err := ParseInt[int](b, math.MinInt, math.MaxInt)
	if err != nil {
		return err
	}
	*f = FlexInt(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexInt8) UnmarshalJSON(b []byte) error {
	value, err := ParseInt[int8](b, math.MinInt8, math.MaxInt8)
	if err != nil {
		return err
	}
	*f = FlexInt8(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexInt16) UnmarshalJSON(b []byte) error {
	value, err := ParseInt[int16](b, math.MinInt16, math.MaxInt16)
	if err != nil {
		return err
	}
	*f = FlexInt16(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexInt32) UnmarshalJSON(b []byte) error {
	value, err := ParseInt[int32](b, math.MinInt32, math.MaxInt32)
	if err != nil {
		return err
	}
	*f = FlexInt32(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexInt64) UnmarshalJSON(b []byte) error {
	value, err := ParseInt[int64](b, math.MinInt64, math.MaxInt64)
	if err != nil {
		return err
	}
	*f = FlexInt64(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexUint) UnmarshalJSON(b []byte) error {
	value, err := ParseUint[uint](b, math.MaxUint)
	if err != nil {
		return err
	}
	*f = FlexUint(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexUint8) UnmarshalJSON(b []byte) error {
	value, err := ParseUint[uint8](b, math.MaxUint8)
	if err != nil {
		return err
	}
	*f = FlexUint8(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexUint16) UnmarshalJSON(b []byte) error {
	value, err := ParseUint[uint16](b, math.MaxUint16)
	if err != nil {
		return err
	}
	*f = FlexUint16(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexUint32) UnmarshalJSON(b []byte) error {
	value, err := ParseUint[uint32](b, math.MaxUint32)
	if err != nil {
		return err
	}
	*f = FlexUint32(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexUint64) UnmarshalJSON(b []byte) error {
	value, err := ParseUint[uint64](b, math.MaxUint64)
	if err != nil {
		return err
	}
	*f = FlexUint64(value)

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexFloat32) UnmarshalJSON(b []byte) error {
	b = ParseBytes(b)
	if len(b) > 0 {
		value, err := strconv.ParseFloat(string(b), 32)
		if err != nil {
			return err
		}
		*f = FlexFloat32(value)
	}

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexFloat64) UnmarshalJSON(b []byte) error {
	b = ParseBytes(b)
	if len(b) > 0 {
		value, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return err
		}
		*f = FlexFloat64(value)
	}

	return nil
}

// UnmarshalJSON 实现了json.Unmarshaler接口
func (f *FlexBool) UnmarshalJSON(b []byte) error {
	b = ParseBytes(b)
	if len(b) > 0 {
		value, _ := strconv.ParseBool(string(b))
		*f = FlexBool(value)
	}

	return nil
}

// UnmarshalJSON 实现了 json.Unmarshaler接口
func (f *FlexObject[T]) UnmarshalJSON(b []byte) error {
	b = ParseBytes(b)
	if len(b) == 0 || string(b) == "null" || string(b) == "[]" || string(b) == "{}" {
		f.ptr = nil
		return nil
	}

	var v T
	if err := Unmarshal(b, &v); err != nil {
		return err
	}
	f.ptr = &v

	return nil
}

// parseInt 解析整数
func ParseInt[T int | int8 | int16 | int32 | int64](b []byte, min, max int64) (T, error) {
	b = ParseBytes(b)
	if len(b) == 0 {
		return 0, nil
	}
	value, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0, err
	}
	if value < min || value > max {
		return 0, &json.UnmarshalTypeError{
			Value: fmt.Sprintf("number %d", value),
			Type:  reflect.TypeOf(T(0)),
		}
	}

	return T(value), nil
}

// parseUint 解析无符号整数
func ParseUint[T uint | uint8 | uint16 | uint32 | uint64](b []byte, max uint64) (T, error) {
	b = ParseBytes(b)
	if len(b) == 0 {
		return 0, nil
	}
	value, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0, err
	}
	if value > max {
		return 0, &json.UnmarshalTypeError{
			Value: fmt.Sprintf("number %d", value),
			Type:  reflect.TypeOf(T(0)),
		}
	}
	return T(value), nil
}

// parseBytes 解析字节数组，去掉首尾空格和引号
func ParseBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	b = bytes.TrimSpace(b)
	if len(b) >= 2 && b[0] == '"' && b[len(b)-1] == '"' {
		return bytes.TrimSpace(b[1 : len(b)-1])
	}

	return b
}

// Unmarshal 将JSON编码的数据解析并存储在v指向的值中。
// 注意：默认情况下该方法将复制给定的缓冲区，如果要更高效地传递JSON，请使用UnmarshalString。
func Unmarshal(data []byte, v any) error {
	return sonic.ConfigStd.Unmarshal(data, v)
}

// UnmarshalString 类似于Unmarshal，不同之处在于data是一个字符串。
func UnmarshalString(data string, v any) error {
	return sonic.ConfigStd.UnmarshalFromString(data, v)
}

// Marshal 返回v的JSON编码字节。
func Marshal(v any) ([]byte, error) {
	return sonic.ConfigStd.Marshal(v)
}

// MarshalString 返回v的JSON编码字符串。
func MarshalString(v any) (string, error) {
	return sonic.ConfigStd.MarshalToString(v)
}

func (f *FlexInt) Value() int {
	return int(*f)
}

func (f *FlexInt8) Value() int8 {
	return int8(*f)
}

func (f *FlexInt16) Value() int16 {
	return int16(*f)
}

func (f *FlexInt32) Value() int32 {
	return int32(*f)
}

func (f *FlexInt64) Value() int64 {
	return int64(*f)
}

func (f *FlexUint) Value() uint {
	return uint(*f)
}

func (f *FlexUint8) Value() uint8 {
	return uint8(*f)
}

func (f *FlexUint16) Value() uint16 {
	return uint16(*f)
}

func (f *FlexUint32) Value() uint32 {
	return uint32(*f)
}

func (f *FlexUint64) Value() uint64 {
	return uint64(*f)
}

func (f *FlexFloat32) Value() float32 {
	return float32(*f)
}

func (f *FlexFloat64) Value() float64 {
	return float64(*f)
}

func (f *FlexBool) Value() bool {
	return bool(*f)
}

func (f *FlexObject[T]) Ptr() *T {
	return f.ptr
}
