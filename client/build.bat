cd ../../../../
set GOPATH=%cd%
cd src/demo/evio_test/client
set GOARCH=amd64
set GOOS=windows
go build -v -ldflags="-s -w"