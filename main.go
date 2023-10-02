package main

func main() {

	setting := openAndRead()

	// NEED_PASSWORD := true
	// PASSWORD := "your_password"
	// port := 3000

	runServer(!setting.NoPassword, setting.Password, int(setting.PortNumber))
}
