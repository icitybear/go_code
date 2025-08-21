# Flexjson

PHP 序列化的 json 字符串经常出现同一个字段，有时是 int，有时是 string。在 golang 这边定义结构体字段类型会很麻烦，使得不得不用 maptostructure 这类的第三方包来处理 json 反序列化。但 maptostructure 的性能很差，在高并发情况下严重影响性能。

```go
str := `{"num":"1"}` // 失败

type A struct {
    Num int `json:"num"`
}

type B struct {
    Num flexjson.FlexInt `json:"num"`
}

var a *A
fmt.Println(json.Unmarshal([]byte(str), &a).Error())
// "json: cannot unmarshal string into Go struct field A.num of type int"

var b *B
json.Unmarshal([]byte(str), &b).Error()
fmt.Println(b.Num.Value() == 1) // true
```

还有一种情况是 PHP 序列化 json 字符串时，可能将 object 序列化为 array，如

```php
$a = [];
if (1 === $b) {
    $a['b'] = $b;
}
echo json_encode($a); // [] or {"b":1}
```

这种情况下，在 golang 中定义为 struct 时，当 json 字符串为 `[]` 时，反序列化会失败。使用 flexjson.FlexObject[T] 即可进行兼容处理。

```go
str := `[]`
type C struct {
    B int `json:"b"`
}

type D struct {
    B flexjson.FlexInt `json:"b"`
}

var c *C
fmt.Println(json.Unmarshal([]byte(str), &c).Error()) // json: cannot unmarshal array into Go value of type C

var d flexjson.FlexObject[D]
json.Unmarshal([]byte(str), &d)
fmt.Println(d.Ptr()) // <nil>
```

该包还集成了 [sonic](https://github.com/bytedance/sonic)，一款由字节跳动开发的高性能 json 序列化、反序列化工具。性能相比golang原生 json 快很多且更节省内存。

```go
str := `{"a":1}`
var a A
err := flexjson.Unmarshal([]byte(str), &a) // 将调用 sonic.Unmarshal 方法
// 或使用
err := flexjson.UnmarshalString(str, &a) // 将调用 sonic.UnmarshalString 方法
```
