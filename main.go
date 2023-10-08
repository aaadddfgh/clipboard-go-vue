package main

import (
	"clipboard-go-vue/src"
	"clipboard-go-vue/src/lib"
)

func main() {

	setting := lib.OpenAndReadSetting()

	// NEED_PASSWORD := true
	// PASSWORD := "your_password"
	// port := 3000

	src.RunServer(!setting.NoPassword, setting.Password, int(setting.PortNumber), setting.Lang)
}
