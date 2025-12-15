package easy_test

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql" // mysql驱动，还有ck驱动
	"gorm.io/gorm"
)

// UserInfo 用户信息
type UserInfo struct {
	ID     uint
	Name   string
	Gender string
	Hobby  string
}

// 连接 插入数据
func TestEasy(t *testing.T) {
	fmt.Println("简单的数据操作")
	source := "root:citybear@(127.0.0.1:13306)/mydata" // 账号：密码 ip端口 数据库名
	dsn := fmt.Sprintf("%s?charset=utf8mb4&readTimeout=%ds&writeTimeout=%ds&parseTime=True&loc=Local", source, 3, 3)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, // 事务
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("数据库资源%#v", db)

	// defer db.Close()

	// 自动迁移
	db.AutoMigrate(&UserInfo{})

	u1 := UserInfo{1, "七米", "男", "篮球"}
	u2 := UserInfo{2, "沙河娜扎", "女", "足球"}
	// 创建记录
	db.Create(&u1)
	db.Create(&u2)
	// 查询
	var u = new(UserInfo) // &UserInfo{}
	db.First(u, 1)        // 根据整型主键查找
	// 找不到记录的情况

	fmt.Printf("%#v\n", u)

	var uu UserInfo
	db.Find(&uu, "hobby=?", "足球")
	fmt.Printf("%#v\n", uu)

	// 绑定模型
	// // 更新
	// db.Model(&u).Update("hobby", "双色球")
	// // 删除
	// db.Delete(&u)
}

// first 和 find 找不到数据 first take last会报错ErrRecordNotFound
func TestFrist(t *testing.T) {
	fmt.Println("简单的数据操作")
	source := "root:citybear@(127.0.0.1:13306)/mydata" // 账号：密码 ip端口 数据库名
	// 读写超时 字符集 以及时区和时间格式自动转换
	dsn := fmt.Sprintf("%s?charset=utf8mb4&readTimeout=%ds&writeTimeout=%ds&parseTime=True&loc=Local", source, 3, 3)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, // 事务
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("数据库资源%#v", db)

	// var uu UserInfo
	// res := db.Find(&uu, "hobby=?", "足球1")
	// fmt.Println(res.Error)
	// fmt.Println(uu)

	// 测试save
	info := &User{}
	db.Find(&info).Where("name = ?", "User1")
	fmt.Println(info)

	updInfo := &User{
		ID:   3,
		Name: "citybear",
		// Status字段 默认0值覆盖 save 可以通过select避免
	}
	cmd := db.Model(&User{}).Save(updInfo)
	fmt.Println(cmd)

	updInfo1 := &User{
		// ID:     8, // 未指定id或默认零值 （主键）的情况下 默认插入 指定的情况下只会更新
		Name:   "citybear2",
		Status: 3,
	}
	cmd2 := db.Model(&User{}).Save(updInfo1)
	fmt.Println(cmd2)
}
