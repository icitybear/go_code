module chsdk

go 1.19

replace sdk => ./sdk

require (
	github.com/satori/go.uuid v1.2.0
	sdk v0.0.0-00010101000000-000000000000
)

require gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
