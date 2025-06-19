package easy_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/hints"
)

type Category struct {
	ID     int
	Code   string `gorm:"column:code"`
	Status int
	// Posts  []Post `gorm:"foreignKey:Code;references:CategoryCode"` // 一对多关系
}

// Posts: define a valid foreign key for relations or implement the Valuer/Scanner interface
func (m *Category) TableName() string {
	return "categorys"
}

type User struct {
	ID     int
	Name   string
	Status int
	Posts  []Post // 一对多关系
}

// 1对多 或者1对1 看字段类型结构体里定义的是 []post还是 post
type Post struct {
	ID           int
	UserID       int
	Title        string
	User         User     `gorm:"foreignKey:UserID"` // 关联到 User 表 references:ID 可以加上 实战一般用指针 *User
	CategoryCode string   `gorm:"column:code"`
	Category     Category `gorm:"foreignKey:CategoryCode;references:Code"` // 关联到 Type 表 references:Code
	PostExt      PostExt  `gorm:"foreignKey:ID;references:PostId"`
}

type PostExt struct {
	ID     int
	Remark string
	PostId int `gorm:"column:post_id"` // 1对1
}

func (m *PostExt) TableName() string {
	return "post_ext"
}

func setupDatabase() (*gorm.DB, error) {

	source := "root:citybear@(127.0.0.1:13306)/mydata" // 账号：密码 ip端口 数据库名
	dsn := fmt.Sprintf("%s?charset=utf8mb4&readTimeout=%ds&writeTimeout=%ds&parseTime=True&loc=Local", source, 3, 3)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, // 事务
		// Logger:                 logger.Default.LogMode(logger.Info), // 全局debug模式
	})
	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{}) // 	"gorm.io/driver/sqlite"
	if err != nil {
		return nil, err
	}

	// db.AutoMigrate(&User{}, &Post{}) // 建表

	// // 填充数据
	// for i := 1; i <= 5; i++ {
	// 	user := User{Name: fmt.Sprintf("User%d", i)}
	// 	db.Create(&user)
	// 	for j := 1; j <= 3; j++ {
	// 		db.Create(&Post{UserID: user.ID, Title: fmt.Sprintf("Post%d-%d", i, j)})
	// 	}
	// }

	// db.AutoMigrate(&Category{})
	// for i := 1; i <= 5; i++ {
	// 	cate := Category{Code: fmt.Sprintf("code%d", i)}
	// 	db.Create(&cate)
	// }
	return db, nil
}

// 不使用 Preload 的查询
func getPostsWithoutPreload(db *gorm.DB) {
	var posts []Post
	// 查询所有帖子
	db.Find(&posts)

	// 打印每个帖子的用户信息
	for _, post := range posts {
		var user User
		// tag: 主动查询
		db.First(&user, post.UserID) // 这会对每个帖子发起一次查询
		fmt.Printf("Post: %s, User: %s\n", post.Title, user.Name)
	}
}

// 使用 Preload 的查询
func getPostsWithPreload(db *gorm.DB) {
	var posts []Post // n:1
	// 使用 Preload 一次性查询所有帖子和用户
	// db.Preload("User").Find(&posts)

	db.Debug().Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("ID", "Name") // 限制获取的字段 不查询 Status
		// 注意，即便只查询 Name 字段，也需要将 ID 字段一并加上。
		// 不加 ID 就会报警告信息： failed to assign association &models.User
	}).Find(&posts)

	// 打印每个帖子的用户信息
	for _, post := range posts {
		fmt.Printf("Post: %s, User: %s, User.Status: %d\n", post.Title, post.User.Name, post.User.Status)
	}
}

func TestMain(t *testing.T) {

	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}

	// fmt.Println("Without Preload:")
	// getPostsWithoutPreload(db) // 不使用 Preload

	fmt.Println("\nWith Preload:")
	getPostsWithPreload(db) // 使用 Preload
}

func getPostsWithPreload2(db *gorm.DB) {
	var posts []Post // n:1
	db.Debug().Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("status = ?", 1) // 这里的字段时 users的status 只获取了 User里status=1的
	}).Find(&posts)
	// SELECT * FROM `users` WHERE status = 1 AND `users`.`id` IN (1,2,3,4,5)
	// SELECT * FROM `posts` // 结果全查了

	// 打印每个帖子的用户信息
	for _, post := range posts {
		fmt.Printf("Post: %s, User: %s, Status: %d \n", post.Title, post.User.Name, post.User.Status)
	}
}
func TestMain2(t *testing.T) {

	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}

	getPostsWithPreload2(db) // 使用 表条件
}

// getPostsWithPreload3  preload  category
func getPostsWithPreload3(db *gorm.DB) {
	var posts []Post // n:1
	// db.Debug().Preload("User", func(tx *gorm.DB) *gorm.DB {
	// 	// 如果写在里面 那么告警无数据 Category: unsupported relations for schema User
	// 	return tx.Where("status = ?", 1).Preload("Category") // 必须User与Category也有关联 status条件时users表的
	// }).Find(&posts)
	// SELECT * FROM `users` WHERE status = 1 AND `users`.`id` IN (1,2,3,4,5)
	// Category: unsupported relations for schema User
	// SELECT * FROM `posts` // 结果全查了 然后关联的数据报错都没了

	db.Debug().Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("status = ?", 1)
	}).Preload("Category", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("status = ?", 1)
	}).Find(&posts) // 写外面 就是Post与Category有关联就行

	// SELECT * FROM `categorys` WHERE status = 1 AND `categorys`.`code` IN ('code1','code2') // 因为posts表code范围code1，code2 零值不会进行预加载
	// SELECT * FROM `users` WHERE status = 1 AND `users`.`id` IN (1,2,3,4,5) // 因为posts表code范围
	// SELECT * FROM `posts` // 结果全查了

	// 打印每个帖子的用户信息
	for _, post := range posts {
		fmt.Printf("Post: %s, User: %s, User.Status: %d, Category: %s, Category.Status: %d \n", post.Title, post.User.Name, post.User.Status, post.Category.Code, post.Category.Status)
	}
}

func TestMain3(t *testing.T) {

	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}

	getPostsWithPreload3(db) // 使用 表条件

}

// 上面都是 n:1 查的Post 预查User
// 下面是 1:1 查的Post 预查PostExt

func getPostsWithPreload4(db *gorm.DB) {
	var posts []Post
	db.Debug().Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("status = ?", 1) // users的status 只获取了 User里status=1的
	}).Preload("PostExt", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("status = ?", 1) // post_ext的status
	}).Find(&posts)
	// 只会预加载 Post的表所有id为条件的post_ext表记录
	// SELECT * FROM `post_ext` WHERE status = 1 AND `post_ext`.`post_id` IN (1,2,3,4,5,6,7,8,9,10,11,12,13,14,15)
	// SELECT * FROM `users` WHERE status = 1 AND `users`.`id` IN (1,2,3,4,5)
	// SELECT * FROM `posts` // 结果全查了

	// 打印每个帖子的用户信息
	for _, post := range posts {
		fmt.Printf("Post: %s, User: %s, User.Status: %d, PostExt: %s \n", post.Title, post.User.Name, post.User.Status, post.PostExt.Remark)
	}
}

func TestMain4(t *testing.T) {

	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}
	getPostsWithPreload4(db)
}

func getPostsWithPreload5(db *gorm.DB) {
	var posts []Post

	// users的status 只获取了 User里status=1的
	db.Debug().Preload("User", "status = ?", 1).
		Preload("PostExt", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("status = ?", 1) // post_ext的status
		}).Find(&posts)
	// SELECT * FROM `post_ext` WHERE status = 1 AND `post_ext`.`post_id` IN (1,2,3,4,5,6,7,8,9,10,11,12,13,14,15)
	// SELECT * FROM `users` WHERE `users`.`id` IN (1,2,3,4,5) AND status = 1
	// SELECT * FROM `posts`

	// 打印每个帖子的用户信息
	for _, post := range posts {
		fmt.Printf("Post: %s, User: %s, User.Status: %d, PostExt: %s \n", post.Title, post.User.Name, post.User.Status, post.PostExt.Remark)
	}
}

func TestMain5(t *testing.T) {

	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}
	getPostsWithPreload5(db)
}

// 使用join进行条件过滤 生成sql 1:1
func getPostsWithJoin(db *gorm.DB) {
	var posts []Post

	db.Debug().Preload("User").Preload("PostExt").
		Joins("left join post_ext t1 on t1.post_id = posts.id").
		Where("t1.status = ?", 1). // 左表字段
		Where("user_id > ?", 3).   // 主表字段
		// Select("posts.*", "t1.remark"). // 可以追加select
		Find(&posts)
	// SELECT * FROM `post_ext` WHERE `post_ext`.`post_id` = 14 // 因为下面的where条件
	// SELECT * FROM `users` WHERE `users`.`id` = 5 // 因为下面的where条件

	// SELECT `posts`.`id`,`posts`.`user_id`,`posts`.`title`,`posts`.`code` FROM `posts`
	// left join post_ext t1 on t1.post_id = posts.id WHERE t1.status = 1 AND user_id > 3

	// 打印每个帖子的用户信息
	for _, post := range posts {
		fmt.Printf("Post: %s, User: %s, User.Status: %d, PostExt: %s \n", post.Title, post.User.Name, post.User.Status, post.PostExt.Remark)
	}
}

func TestMain6(t *testing.T) {

	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}
	getPostsWithJoin(db)
}

func TestSql(t *testing.T) {

	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}
	// 方案1: Debug模式 全局和临时
	var posts []Post
	// 方案2:  必须在tryRun模式 默认情况 GORM 会在执行 SQL 后清空 Statement.SQL，导致后续访问为空。
	result := db.Session(&gorm.Session{DryRun: true}).Preload("User").Preload("PostExt").
		Joins("left join post_ext t1 on t1.post_id = posts.id").
		Where("t1.status = ?", 1). // 左表字段
		Where("user_id > ?", 3).   // 主表字段
		Find(&posts)

	sql := result.Statement.SQL.String()
	vars := result.Statement.Vars
	fmt.Println(sql)
	spew.Println(vars)
	// SELECT `posts`.`id`,`posts`.`user_id`,`posts`.`title`,`posts`.`code` FROM `posts` left join post_ext t1 on t1.post_id = posts.id WHERE t1.status = ? AND user_id > ?
	// [1 3]
	realSql := db.Dialector.Explain(sql, vars...)
	fmt.Println(realSql)

	spew.Println(posts) // tryRun模式 不会去真正执行

	// 方案3: 无需 DryRun，但仅支持查询类操作（如 Find、First)
	// ToSQL 通过闭包接收 GORM 的查询链式调用，并返回最终构建的 SQL 字符串
	sql2 := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		// 这里真正执行了
		var total int64
		res := tx.Table("categorys").Where("1=1").Count(&total)
		fmt.Println(total) // 不会真正执行
		return res
	})
	fmt.Println(sql2)

}

func TestHint(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}
	var uu Category
	// .Hint undefined (type *gorm.DB has no field or method Hint //
	// res := db.Hint("FORCE INDEX(idx_code_status)").Find(&uu, "code=?", "code2") // 报错
	// res := db.Find(&uu, "code=?", "code2")
	// https://learnku.com/docs/gorm/v2/hints/9769
	txDb := db.Debug().Clauses(hints.ForceIndex("idx_code_status")).Where("code=?", "code2")
	var total int64
	err = txDb.Table("categorys").Count(&total).Error
	fmt.Println(err)
	fmt.Println(total)
	// .Clauses(hints.ForceIndex("idx_code_status")) txDb会继承
	res := txDb.Find(&uu) // 尽量不使用内嵌条件
	fmt.Println(res.Error)
	fmt.Println(uu)

	// 未报错 正确
	var ct []Category
	db.Debug().Raw("SELECT * FROM categorys FORCE INDEX (idx_code_status) WHERE code = ?", "code2").Scan(&ct)
	fmt.Println(ct)
	// scan要显示指定表
	var total2 sql.NullFloat64 // 命名不需要与as字段一致
	err = db.Debug().Table("categorys").Select("sum(status) as total").Where("id > ?", 0).Scan(&total2).Error
	if err != nil {
		fmt.Println(err)
	}
	if total2.Valid {
		fmt.Println(float32(total2.Float64))
	} else {
		fmt.Println("end")
	}

}
