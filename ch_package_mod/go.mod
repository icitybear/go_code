module packmod

go 1.17

require test_app v1.0.0

require (
	github.com/satori/go.uuid v1.2.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
///mnt/c/myapps/go/go_code/code/ch_package_mod/custom/test_app //相对路径 绝对路径也可以
replace test_app => ./custom/test_app
