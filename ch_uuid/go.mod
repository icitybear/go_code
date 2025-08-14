module ch_uuid

go 1.19

require (
	github.com/google/uuid v1.6.0
	github.com/rs/xid v1.6.0
	github.com/satori/go.uuid v1.2.0
	kid v1.0.0
)

require gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect

replace kid => ./kid // 相对路径 绝对路径也可以
