package lib

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
)

type Setting struct {
	Password   string `json:"Password" validate:"require_unless=NoPassword true, max=20,min=6"`
	NoPassword bool   `json:"NoPassword"`
	PortNumber uint16 `json:"PortNumber"`
}

func OpenAndReadSetting() Setting {

	file, err := os.Open("./config.json")
	if err != nil {
		log.SetPrefix("[App] [WARNNING] ")
		log.Println("No config file, will run without password on port 3000")

		return Setting{
			Password:   "",
			NoPassword: true,
			PortNumber: 3000,
		}
	}
	defer file.Close()

	jsonDcoder := json.NewDecoder(file)

	var setting Setting
	err = jsonDcoder.Decode(&setting)

	if _, syntaxError := err.(*json.SyntaxError); syntaxError {
		log.SetPrefix("[App] [WARNNING] ")
		log.Println("config file error")
		log.Println("No config file, will run without password on port 3000")
		return Setting{
			Password:   "",
			NoPassword: true,
			PortNumber: 3000,
		}
	}

	validate := validator.New()
	if setting.NoPassword {
		return setting
	}
	if err := validate.Struct(setting); err != nil {
		panic(err.Error())
	}
	return setting
}
