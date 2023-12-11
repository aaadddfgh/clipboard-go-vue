set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=0

go build -o ".\release\linux-arm64\clipboard.out" 

copy .\config.json ".\release\linux-arm64\config.json"

xcopy .\public ".\release\linux-arm64\public\" /-Y