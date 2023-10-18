package main

import (
	"os/user"
	"path"
)

func GetDefaultPath() string {
	u, _ := user.Current()
	return path.Join(u.HomeDir, "Documents")
}
