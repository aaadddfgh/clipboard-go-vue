set GOOS=windows
set GOARCH=386
set CGO_ENABLED=0

go build -o ".\release\x86-64\clipboard.exe"

copy .\config.json ".\release\x86-64\config.json"

xcopy ".\public\" ".\release\x86-64\public\" /-Y