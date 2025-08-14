# kid

基于 k8s pod-name 规则的本地生成分布式唯一ID算法

本项目参考 [雪花ID](https://en.wikipedia.org/wiki/Snowflake_ID) 和 [MangoDB ObjectID](https://www.mongodb.com/docs/manual/reference/bson-types/#std-label-objectid)

雪花ID的优点在于可使用 int64 进行存储，可在数据库中作为自增ID使用，但机器ID只有1024个，不依赖外部系统的话，机器ID容易冲突。

MangoDB ObjectID 使用96个比特，只能用字符串(char, varchar) 或 二进制(blob) 存储，但机器ID范围比雪花ID大很多，避免了机器ID冲突。

参考 k8s pod-name 的生成规则后，本项目使用 pod-name 后面的[5个随机字符](https://github.com/kubernetes/kubernetes/blob/71b8ad965ef0e4cc2a2c22a68cf7db8d1f8c6dc8/staging/src/k8s.io/apimachinery/pkg/util/rand/rand.go#L83)作为机器ID，保证在同一服务集群下，机器ID不冲突。

输出时使用自定义的 base32 编码(默认编码字典`sh2tqf6lownky4pcg8xvze5jma9rd7bu`)，输出为 13 ～ 16 个字符的字符串。

## 数据格式

```text
┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
│ 0 │ 1 │ 2 │...│40 │41 │42 │43 │44 │...│64 │65 │66 │67 │68 │...│77 │78 │79 │
└───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
  └───────────────────┘   └───────────────────┘   └───────────────────────┘
    毫秒级时间戳 (42位)          机器ID (24位)              自增ID (14位)
        位 0-41                  位 42-65                  位 66-79
```

| 字段 | 位范围 | 位数 | 数值范围 | 说明 |
|------|--------|------|----------|------|
| 毫秒级时间戳 | 0-41 | 42位 | 0 ~ 4,398,046,511,103 | 约139年的时间范围 |
| 机器ID | 42-65 | 24位 | 0 ~ 16,777,215 | 支持约1677万台机器 |
| 自增ID | 66-79 | 14位 | 0 ~ 16,383 | 每毫秒可生成16384个唯一ID<br>自增ID和MangoDB ObjectID一样，单调递增，不随时间戳进行重置（雪花ID的行为） |

> 同一秒内可生成 16,777,216 个ID，在同一毫秒内可生成 16,384 个ID。

## 使用

```go
import "git.internal.taqu.cn/go-modules/kid"

uuid := kid.New().String()
fmt.Println(uuid) // gspr6qygvhssh
```

设置编码字典

```go
// 默认编码字典 sh2tqf6lownky4pcg8xvze5jma9rd7bu
err := kid.SetEncodingDict("0123456789abcdefghijklmnopqrstuv")
if err != nil {
    panic(fmt.Sprintf("设置kid编码字典失败. err:", err))
}
```

设置纪元，默认时间纪元为 1752741566491(2025-07-17T16:39:26.491+0800)，上面的数据格式里说前面42位时间戳最多能用到 2109-05-15T15:35:11+0800，扣去默认的时间纪元后，实际可用到 2164-11-29T00:14:37+0800。如果不想用默认的时间纪元，可以通过以下方式设置

```go
// 不使用时间纪元
kid.Eopch = 0

id := kid.New()
fmt.Println(id) // gspr6qygvhssh
```

设置机器ID，如果希望通过环境变量设置机器ID，可通过下面方式设置

```go
// 假设环境变量名为 KID_MACHINE_ID
mid := os.Getenv("KID_MACHINE_ID")
if mid == "" {
    return nil
}

num, err := strconv.Atoi(mid)
if err != nil {
    panic("KID_MACHINE_ID value is set to not a number")
}
kid.SetMachineID(num)
id := kid.New()
fmt.Println(id)
```
