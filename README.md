# clipboard-go-vue

A web clipboard 

## Features
### secure communication
This clipboard protect your communication even you don't use HTTPS/SSL to protect your infomation.
(can't prevent Man-in-the-MiddleAttack)  
And protect your clipboard with password (alternative).  

### config.json
In config.json
``` json
{
    "Password": "your_password",
	"PortNumber": 3000,
}
```

``` go
type Setting struct {
	Password   string `json:"Password" validate:"required_unless=NoPassword true,max=20,min=6"`
	NoPassword bool   `json:"NoPassword"`
	PortNumber uint16 `json:"PortNumber"`
	Lang       string `json:"Lang"`
}
```

## Build
If you don't want to build yourself, you can download flies in release.  
Else:
### Requirement
1. node >= 14
2. go >= 1.20

### Steps  
1. build html  
Pull branch `vue`, in the dir you pulled `vue` branch, run
```
npm install
npm bulid
```
after that, html and other files should be in `./dist` directory.   

2. build server

Pull branch `main`, in the directory you pulled `main` branch, run
```
go install
go build
```
you should see `clipboard.exe`.   

3. run

Move files you got in 1st step to `./public` in same directory as `clipboard.exe`, and run it.


