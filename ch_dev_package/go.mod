module ch_dev_package

go 1.23.6

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/google/uuid v1.6.0
	github.com/rs/xid v1.6.0
	github.com/satori/go.uuid v1.2.0
	github.com/spaolacci/murmur3 v1.1.0
	im-robot v0.0.0-00010101000000-000000000000
	kid v1.0.0
)

require (
	github.com/go-resty/resty/v2 v2.16.5 // indirect
	golang.org/x/net v0.33.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace kid => ./kid // 相对路径 绝对路径也可以

replace im-robot => ./im-robot
