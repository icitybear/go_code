package flexjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\t' || c == '\r' || c == '\n')
}

func noSpace(c rune) rune {
	if isSpace(byte(c)) { // only used for ascii
		return -1
	}
	return c
}

type unmarshalTest struct {
	in  string
	ptr interface{}
	out interface{}
	err error
}

type zeroType struct {
	Num   FlexInt
	UNum  FlexUint
	Float FlexFloat64
	Bool  FlexBool
	Omit  FlexInt `json:"omit,omitempty"`
}

var unmarshalTests = []unmarshalTest{
	{`{"num":"123"}`, new(map[string]FlexInt), map[string]FlexInt{"num": 123}, nil},
	{`{"num":123}`, new(map[string]FlexInt), map[string]FlexInt{"num": 123}, nil},
	{`{"num":"12"}`, new(map[string]FlexInt8), map[string]FlexInt8{"num": 12}, nil},
	{`{"num":12}`, new(map[string]FlexInt8), map[string]FlexInt8{"num": 12}, nil},
	{`{"num":"23134"}`, new(map[string]FlexInt16), map[string]FlexInt16{"num": 23134}, nil},
	{`{"num":23134}`, new(map[string]FlexInt16), map[string]FlexInt16{"num": 23134}, nil},
	{`{"num":"23134"}`, new(map[string]FlexInt32), map[string]FlexInt32{"num": 23134}, nil},
	{`{"num":23134}`, new(map[string]FlexInt32), map[string]FlexInt32{"num": 23134}, nil},
	{`{"num":"-23134"}`, new(map[string]FlexInt64), map[string]FlexInt64{"num": -23134}, nil},
	{`{"num":23134}`, new(map[string]FlexInt64), map[string]FlexInt64{"num": 23134}, nil},
	{`{"num":"123"}`, new(map[string]FlexUint), map[string]FlexUint{"num": 123}, nil},
	{`{"num":123}`, new(map[string]FlexUint), map[string]FlexUint{"num": 123}, nil},
	{`{"num":"12"}`, new(map[string]FlexUint8), map[string]FlexUint8{"num": 12}, nil},
	{`{"num":12}`, new(map[string]FlexUint8), map[string]FlexUint8{"num": 12}, nil},
	{`{"num":"23134"}`, new(map[string]FlexUint16), map[string]FlexUint16{"num": 23134}, nil},
	{`{"num":23134}`, new(map[string]FlexUint16), map[string]FlexUint16{"num": 23134}, nil},
	{`{"num":"23134"}`, new(map[string]FlexUint32), map[string]FlexUint32{"num": 23134}, nil},
	{`{"num":23134}`, new(map[string]FlexUint32), map[string]FlexUint32{"num": 23134}, nil},
	{`{"num":"23134"}`, new(map[string]FlexUint64), map[string]FlexUint64{"num": 23134}, nil},
	{`{"num":23134}`, new(map[string]FlexUint64), map[string]FlexUint64{"num": 23134}, nil},
	{`{"num":3.1415926535}`, new(map[string]FlexFloat32), map[string]FlexFloat32{"num": 3.1415926535}, nil},
	{`{"num":"3.1415926535"}`, new(map[string]FlexFloat32), map[string]FlexFloat32{"num": 3.1415926535}, nil},
	{`{"num":3.1415926535}`, new(map[string]FlexFloat64), map[string]FlexFloat64{"num": 3.1415926535}, nil},
	{`{"num":"3.1415926535"}`, new(map[string]FlexFloat64), map[string]FlexFloat64{"num": 3.1415926535}, nil},
	{`{"bool":true,"str":"True","sstr":"t","snum":"1","num":1}`, new(map[string]FlexBool), map[string]FlexBool{"bool": true, "str": true, "sstr": true, "snum": true, "num": true}, nil},
	{`{"bool":false,"str":"False","sstr":"f","snum":"0","num":0,"any_else":"ksdf"}`, new(map[string]FlexBool), map[string]FlexBool{"bool": false, "str": false, "sstr": false, "snum": false, "num": false, "any_else": false}, nil},

	// 0值
	{`{"num":"","unum":"","float":"","bool":""}`, new(zeroType), zeroType{Num: 0, UNum: 0, Float: 0, Bool: false, Omit: 0}, nil},

	// 超出表述范围
	{`{"num":"257"}`, new(map[string]FlexUint8), nil, fmt.Errorf("json: cannot unmarshal number 257 into Go value of type uint8")},
	{`{"num":"-129"}`, new(map[string]FlexInt8), nil, fmt.Errorf("json: cannot unmarshal number -129 into Go value of type int8")},

	// parseInt 错误
	{`{"num":"abc"}`, new(map[string]FlexInt), nil, fmt.Errorf(`strconv.ParseInt: parsing "abc": invalid syntax`)},
	{`{"num":"abc"}`, new(map[string]FlexInt8), nil, fmt.Errorf(`strconv.ParseInt: parsing "abc": invalid syntax`)},
	{`{"num":"abc"}`, new(map[string]FlexInt16), nil, fmt.Errorf(`strconv.ParseInt: parsing "abc": invalid syntax`)},
	{`{"num":"abc"}`, new(map[string]FlexInt32), nil, fmt.Errorf(`strconv.ParseInt: parsing "abc": invalid syntax`)},
	{`{"num":"abc"}`, new(map[string]FlexInt64), nil, fmt.Errorf(`strconv.ParseInt: parsing "abc": invalid syntax`)},
	{`{"num":"-2"}`, new(map[string]FlexUint), nil, fmt.Errorf(`strconv.ParseUint: parsing "-2": invalid syntax`)},
	{`{"num":"-2"}`, new(map[string]FlexUint8), nil, fmt.Errorf(`strconv.ParseUint: parsing "-2": invalid syntax`)},
	{`{"num":"-2"}`, new(map[string]FlexUint16), nil, fmt.Errorf(`strconv.ParseUint: parsing "-2": invalid syntax`)},
	{`{"num":"-2"}`, new(map[string]FlexUint32), nil, fmt.Errorf(`strconv.ParseUint: parsing "-2": invalid syntax`)},
	{`{"num":"-2"}`, new(map[string]FlexUint64), nil, fmt.Errorf(`strconv.ParseUint: parsing "-2": invalid syntax`)},

	// parseFloat 错误
	{`{"num":"abc"}`, new(map[string]FlexFloat32), nil, fmt.Errorf(`strconv.ParseFloat: parsing "abc": invalid syntax`)},
	{`{"num":"xyz"}`, new(map[string]FlexFloat64), nil, fmt.Errorf(`strconv.ParseFloat: parsing "xyz": invalid syntax`)},
}

func TestUnmarshal(t *testing.T) {
	for i, tt := range unmarshalTests {
		// check decode result for invalid jsons
		if !json.Valid([]byte(tt.in)) || tt.ptr == nil {
			var sv interface{}
			err := UnmarshalString(tt.in, &sv)
			if err == nil && tt.err != nil {
				t.Errorf("test json #%d: %v, %v, want %v", i, tt.in, err, tt.err)
			}
			continue
		}
		typ := reflect.TypeOf(tt.ptr)
		if typ.Kind() != reflect.Ptr {
			t.Errorf("#%d: unmarshalTest.ptr %T is not a pointer type", i, tt.ptr)
			continue
		}
		typ = typ.Elem()

		v := reflect.New(typ)
		if !reflect.DeepEqual(tt.ptr, v.Interface()) {
			// There's no reason for ptr to point to non-zero data,
			// as we decode into new(right-type), so the data is
			// discarded.
			// This can easily mean tests that silently don't test
			// what they should. To test decoding into existing
			// data, see TestPrefilled.
			t.Errorf("#%d: unmarshalTest.ptr %#v is not a pointer to a zero value", i, tt.ptr)
			continue
		}

		if err := UnmarshalString(tt.in, v.Interface()); err != nil || tt.err != nil {
			if err != nil && tt.err != nil && err.Error() == tt.err.Error() {
				continue
			}

			spew.Dump(tt)
			t.Fatalf("#%d: %v, want %v", i, err, tt.err)
			continue
		}
		if !reflect.DeepEqual(v.Elem().Interface(), tt.out) {
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v.Elem().Interface(), tt.out)
			data, _ := MarshalString(v.Elem().Interface())
			println(data)
			data, _ = MarshalString(tt.out)
			println(data)
			continue
		}

		// Check round trip also decodes correctly.
		if tt.err == nil {
			enc, err := Marshal(v.Interface())
			if err != nil {
				t.Errorf("#%d: error re-marshaling: %v", i, err)
				continue
			}
			vv := reflect.New(reflect.TypeOf(tt.ptr).Elem())
			if err := Unmarshal(enc, vv.Interface()); err != nil {
				t.Errorf("#%d: error re-unmarshaling %#q: %v", i, enc, err)
				continue
			}
			if !reflect.DeepEqual(v.Elem().Interface(), vv.Elem().Interface()) {
				t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v.Elem().Interface(), vv.Elem().Interface())
				t.Errorf("     In: %q", strings.Map(noSpace, tt.in))
				t.Errorf("Marshal: %q", strings.Map(noSpace, string(enc)))
				continue
			}

			str, err := MarshalString(v.Interface())
			if err != nil {
				t.Errorf("#%d: error re-marshaling: %v", i, err)
				continue
			}
			vv = reflect.New(reflect.TypeOf(tt.ptr).Elem())
			if err := UnmarshalString(str, vv.Interface()); err != nil {
				t.Errorf("#%d: error re-unmarshaling %#q: %v", i, str, err)
				continue
			}
			if !reflect.DeepEqual(v.Elem().Interface(), vv.Elem().Interface()) {
				t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v.Elem().Interface(), vv.Elem().Interface())
				t.Errorf("     In: %q", strings.Map(noSpace, tt.in))
				t.Errorf("Marshal: %q", strings.Map(noSpace, str))
				continue
			}
		}
	}
}

func TestEmptyBytes(t *testing.T) {
	var x FlexInt
	if err := x.UnmarshalJSON([]byte{}); err != nil {
		t.Errorf("UnmarshalJSON([]byte{}): %v", err)
	} else if x != 0 {
		t.Errorf("UnmarshalJSON([]byte{}) = %d; want 0", x)
	}
}

func TestEmptyObject(t *testing.T) {
	type myStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type emptyObject struct {
		Cate   string               `json:"cate"`
		Object FlexObject[myStruct] `json:"object"`
	}

	str := `{"cate":"test","object":[]}`
	var out emptyObject
	if err := Unmarshal([]byte(str), &out); err != nil {
		t.Errorf("解析空对象失败: %v, in:%s", err, str)
	}
	if out.Object.Ptr() != nil {
		t.Errorf("解析空对象失败，非预期的output: %v", out)
	}
	t.Logf("%#v", out.Object)

	str = `{"cate":"ok","object":{"name":"test"}}`
	var out2 emptyObject
	if err := Unmarshal([]byte(str), &out2); err != nil {
		t.Errorf("解析非空对象失败: %v, in:%s", err, str)
	}
	if out2.Object.Ptr() == nil {
		t.Errorf("解析非空对象失败，非预期的output: %v", out2)
	}
	if out2.Object.Ptr().Name != "test" {
		t.Errorf("解析非空对象失败, 非预期的output: %v", out2.Object)
	}
	t.Logf("%#v", *out2.Object.ptr)

	str = `{"cate":"test","object":{}}`
	var out3 emptyObject
	if err := Unmarshal([]byte(str), &out3); err != nil {
		t.Errorf("解析空对象失败: %v, in:%s", err, str)
	}
	if out3.Object.Ptr() != nil {
		t.Errorf("解析空对象失败，非预期的output: %v", out3)
	}
	t.Logf("%#v", out3.Object)

	str = `{"cate":"test","object":[]}`
	var out4 emptyObject
	if err := Unmarshal([]byte(str), &out4); err != nil {
		t.Errorf("解析空对象失败: %v, in:%s", err, str)
	}
	if out4.Object.Ptr() != nil {
		t.Errorf("解析空对象失败，非预期的output: %v", out4)
	}
	t.Logf("%#v", out4.Object)
}
