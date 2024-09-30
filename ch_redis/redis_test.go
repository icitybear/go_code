package easy_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
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
	// isSuccess, err := rdb.SetNX(ctx, redisKey, "csx", 10*time.Second).Result()  bool, error
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
