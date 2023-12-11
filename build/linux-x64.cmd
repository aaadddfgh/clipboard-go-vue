set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

go build -o ".\release\linux-x64\clipboard.out" 

copy .\config.json ".\release\linux-x64\config.json"

xcopy .\public ".\release\linux-x64\public\" /Y /S /I