package easy_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/redis/go-redis/v9"
)

// 连接 插入数据
func TestRedis(t *testing.T) {
	fmt.Println("简单的缓存操作")
	rdb := redis.NewClient(&redis.Options{
		Addr:         "r-8vbm8ugjj62s3kjhdq.redis.zhangbei.rds.aliyuncs.com:6379",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	})
	fmt.Printf("缓存资源%#v", rdb)
	ctx := context.Background()
	redisKey := "citybear:test"

	isSuccess, err := rdb.Set(ctx, redisKey, "ccccc", 10*time.Second).Result() // string, error

	str := fmt.Sprintf("[SetNX] %s, %+v, %+v", redisKey, isSuccess, err)
	fmt.Println(str)
	if err != nil {
		fmt.Println(err.Error())
	}

	intCmd := rdb.Get(ctx, redisKey)
	if intCmd.Err() != nil {
		fmt.Println(intCmd.Err())
	}
	fmt.Println(intCmd.Val())
	td := rdb.TTL(ctx, redisKey)
	fmt.Println(td.Val())
	// hash结构
	redisKey2 := "hash:test"
	// []string{"name", "张三", "age", "18", "score", "99.06"}
	// map[string]interface{}{"key1": "value1", "key2": "value2"})
	setRes, err := rdb.HSet(ctx, redisKey2, map[string]interface{}{"name": "张三", "age": 18, "score": 99.06}).Result()
	str2 := fmt.Sprintf("[HSet] %s, %+v, %+v", redisKey2, setRes, err)
	fmt.Println(str2)
	rdb.Expire(ctx, redisKey2, 30*time.Second)
	td2 := rdb.TTL(ctx, redisKey2)
	fmt.Println(td2.Val())

	s1 := rdb.HGet(ctx, redisKey2, "score")
	fmt.Println(s1.Val())

	s2 := rdb.HMGet(ctx, redisKey2, "name", "age")
	fmt.Println(s2.Val()) // [张三 18]

	s3 := rdb.HGetAll(ctx, redisKey2)
	fmt.Println(s3.Val()) // map[string]string
	// map可以跟struct转化 直接获取

	// 如果指定获取几个 就封装下
	s4, err := MyHMGet(ctx, rdb, redisKey2, "name", "age")
	fmt.Println(s4, err)
}

func MyHMGet(ctx context.Context, rdb *redis.Client, key string, fields ...string) (map[string]interface{}, error) {
	s := rdb.HMGet(ctx, key, fields...)
	if s.Err() != nil {
		return nil, s.Err()
	}
	slicesArr := s.Val()
	res := map[string]interface{}{}
	for pos, val := range fields {
		res[val] = slicesArr[pos]
	}
	return res, nil
}

func TestRedisTtl(t *testing.T) {

	rdb := redis.NewClient(&redis.Options{
		Addr:         "r-8vbm8ugjj62s3kjhdq.redis.zhangbei.rds.aliyuncs.com:6379",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	})

	ctx := context.Background()
	redisKey := "v3.0/citybear/uri1"                                           // tag: key格式 规范基本是用:隔开 复杂结构体可以序列化为字符串作为key
	isSuccess, err := rdb.Set(ctx, redisKey, "ccccc", 10*time.Second).Result() // string, error

	str := fmt.Sprintf("[SetNX] %s, %+v, %+v", redisKey, isSuccess, err)
	fmt.Printf("日志%s \n", str)
	if err != nil {
		spew.Println(err)
		return
	}
	testKey := "citybear:test"
	ttlCmd := rdb.TTL(ctx, "ab")
	err = ttlCmd.Err()
	val := ttlCmd.Val()
	// 不存在-2 永不过期-1
	fmt.Printf("不存在的情况 ttlCmd:%+v, err:%v, val:%v, %v \n", ttlCmd, err, val, val.Nanoseconds()) // -2ns, -2

	rdb.SetNX(ctx, testKey, "ceshi", 0).Result() // 永不过期时间小于0
	ttlCmd1 := rdb.TTL(ctx, testKey)
	err = ttlCmd1.Err()
	val1 := ttlCmd1.Val()
	fmt.Printf("永不过期的情况 ttlCmd1:%+v, err:%v, val:%v, %v \n", ttlCmd1, err, val1, val1.Nanoseconds()) // -1ns, -1

	ttlCmd2 := rdb.TTL(ctx, redisKey)
	err = ttlCmd2.Err()
	val2 := ttlCmd2.Val()
	fmt.Printf("有效期的情况 ttlCmd2:%+v, err:%v, val:%v, %v \n", ttlCmd2, err, val2, val2.Nanoseconds()) // 10s, 10000000000

	// 时间戳参数范围
	rdb.SetEx(ctx, testKey, "ceshi3", time.Second*time.Duration(100)).Result() // 100s
	ttlCmd3 := rdb.TTL(ctx, testKey)
	err = ttlCmd3.Err()
	val3 := ttlCmd3.Val()
	fmt.Printf("时间戳参数范围 ttlCmd3:%+v, err:%v, val:%v, %v \n", ttlCmd3, err, val3, val3.Nanoseconds()) // val:1m40s, 100000000000

	rdb.SetEx(ctx, testKey, "ceshi4", time.Second*time.Duration(1762833298)).Result() // 不管大于还是小于当前时间戳 这种情况就是当作正常数值
	ttlCmd4 := rdb.TTL(ctx, testKey)
	err = ttlCmd4.Err()
	val4 := ttlCmd4.Val()
	ts := val4 / time.Second
	// val:489675h54m58s, 1762833298000000000, 秒级:1762833298
	fmt.Printf("时间戳参数范围 ttlCmd4:%+v, err:%v, val:%v, %v, 秒级:%d \n", ttlCmd4, err, val4, val4.Nanoseconds(), ts)
	// tag:如果要明确的话 就是自己计算保存的时间秒数大小
}

func TestRedisLua(t *testing.T) {

	rdb := redis.NewClient(&redis.Options{
		Addr:         "r-8vbm8ugjj62s3kjhdq.redis.zhangbei.rds.aliyuncs.com:6379",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	})
	var err error
	ctx := context.Background()
	// 设置缓存
	lockCommand := `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
end
return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])`
	testKey := "citybear:test"

	// tag:相当于del 必须有return 不然会都是 err:redis: nil, val:<nil> 改之后 err:<nil>, val:0 // del操作返回 0,1
	delCmd := rdb.Eval(ctx, `return redis.call("DEL", KEYS[1])`, []string{testKey}, []string{})
	fmt.Printf("delCmd err:%+v, val:%v \n", delCmd.Err(), delCmd.Val())

	// setex单位是秒 不管时间戳是大于还是小于 err:<nil>, val:OK    setex key ttl val
	setexCmd := rdb.Eval(ctx, `local new_tokens = 1762876800 return redis.call("setex", KEYS[1], ARGV[1], new_tokens)`, []string{testKey}, []string{"100"})
	err = setexCmd.Err()
	fmt.Printf("setexCmd err:%+v, val:%v \n", setexCmd.Err(), setexCmd.Val())

	// 执行命令
	resCmd := rdb.Eval(ctx, lockCommand,
		[]string{testKey},            // key数组
		[]string{"ceshi2", "100000"}) // 1000*10 相当于10s

	err = resCmd.Err()
	// 已经存在 没设置成功就是  err:redis: nil, val:<nil> 此时resCmd是整个lua脚本
	// 设置成功 err:<nil>, val:OK set操作返回OK
	fmt.Printf("resCmd err:%v, val:%v \n", resCmd.Err(), resCmd.Val())
	if errors.Is(err, redis.Nil) {
		spew.Printf("已经存在 没设置成功") // 此时错误就为redis.Nil
		return
	}
	// 上下文超时取消等
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		spew.Printf("上下问超时")
		return
	}
	// 执行成功就是 err:<nil>, val:OK
	if err != nil {
		spew.Printf("eval执行失败")
		return
	}

	str, _ := rdb.Get(ctx, testKey).Result()
	fmt.Println(str)
	ttlCmd3 := rdb.TTL(ctx, testKey)
	err = ttlCmd3.Err()
	val3 := ttlCmd3.Val()
	fmt.Printf("ttlCmd3:%+v, err:%v, val:%v, %v \n", ttlCmd3, err, val3, val3.Nanoseconds())
	_ = lockCommand
	// 删除缓存 数值1==数值2 and xxx【0,1】 or 0
	// delCommand := `return redis.call("GET", KEYS[1]) == ARGV[1] and redis.call("DEL", KEYS[1]) or 0`
}

func TestRedisLua2(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "r-8vbm8ugjj62s3kjhdq.redis.zhangbei.rds.aliyuncs.com:6379",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	})
	var err error
	ctx := context.Background()

	test1Key := "{citybear:test}"
	test2Key := "{citybear:test}:ts"

	getCmd := rdb.Get(ctx, test1Key)
	// getCmd := rdb.Eval(ctx, `return redis.call("get", KEYS[1])`, []string{test1Key}, []string{})
	fmt.Printf("getCmd err:%v, val:%v \n", getCmd.Err(), getCmd.Val())
	ttlCmd1 := rdb.TTL(ctx, test1Key)
	err = ttlCmd1.Err()
	val1 := ttlCmd1.Val()
	fmt.Printf("ttlCmd1 err:%v, val:%v, %v \n", err, val1, val1.Nanoseconds())

	resCmd := rdb.Eval(ctx,
		scriptCommand,
		[]string{
			test1Key,
			test2Key,
		},
		// 参数
		[]string{
			"10",                                     // 速率 qps
			"1",                                      // 容量     容量-请求数 就是要存的
			strconv.FormatInt(time.Now().Unix(), 10), // 时间戳 秒 记录刷新的时间值
			"1",                                      // 每次的请求数
		})
	err = resCmd.Err()
	fmt.Printf("resCmd err:%v, val:%v \n", resCmd.Err(), resCmd.Val())
	if errors.Is(err, redis.Nil) {
		spew.Printf("redis nil") // 此时错误就为redis.Nil
		return
	}
	// 上下文超时取消等
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		spew.Printf("上下问超时")
		return
	}
	// 执行成功就是 err:<nil>, val:OK
	if err != nil {
		spew.Printf("eval执行失败")
		return
	}
	str2, _ := rdb.Get(ctx, test1Key).Result()
	fmt.Println(str2)
	ttlCmd2 := rdb.TTL(ctx, test1Key)
	err = ttlCmd2.Err()
	val2 := ttlCmd2.Val()
	fmt.Printf("ttlCmd2 err:%v, val:%v, %v \n", err, val2, val2.Nanoseconds())
}

// math.ceil(fill_time*2)
var scriptCommand = `local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local fill_time = capacity/rate
local ttl = 20

local last_tokens = tonumber(redis.call("get", KEYS[1]))
if last_tokens == nil then
    last_tokens = capacity
end

local last_refreshed = tonumber(redis.call("get", KEYS[2]))
if last_refreshed == nil then
    last_refreshed = 0
end

local delta = math.max(0, now-last_refreshed)
local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
local allowed = filled_tokens >= requested
local new_tokens = filled_tokens
if allowed then
    new_tokens = filled_tokens - requested
end

redis.call("setex", KEYS[1], ttl, new_tokens)

redis.call("setex", KEYS[2], ttl, now)

return allowed`
