module wenjian

go 1.19

require (
	app_sdk v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.28.0
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace app_sdk => ./app_sdk/
