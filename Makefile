
LINUX_BIN=app-updater
LINUX_SERVICE=my-service
LINUX_APP=my-app

WIN_BIN=app-updater.exe
WIN_SERVICE=my-service
WIN_APP=My Application.exe

#main output directory
OUT=out

VERSION=$(shell git describe --tags --always --long --dirty)
#VERSION=`git describe --tags --long`

# MAKEFLAGS += --silent

# Builds the project
all:	build-linux build-win64 build-win32

build-linux: 
	echo "building for linux/amd64"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(OUT)/linux/amd64/$(LINUX_BIN) -ldflags '-w -s -X "main.Version=$(VERSION)" -X "main.runtimeApp=$(LINUX_APP)" -X "main.runtimeService=$(LINUX_SERVICE)"' main.go

build-win64:
	echo "building for windows/amd64"
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o $(OUT)/windows/amd64/$(WIN_BIN) -ldflags '-w -s -X "main.Version=$(VERSION)" -X "main.runtimeApp=$(WIN_APP)" -X "main.runtimeService=$(WIN_SERVICE)"' main.go
	
build-win32:
	echo "building for windows/386"
	env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a -o $(OUT)/windows/386/$(WIN_BIN) -ldflags '-w -s -X "main.Version=$(VERSION)" -X "main.runtimeApp=$(WIN_APP)" -X "main.runtimeService=$(WIN_SERVICE)"' main.go



# Cleans our project: deletes binaries
clean:
	if [ -d $(OUT) ];then echo "removing directory $(OUT)"; rm -rf $(OUT) ;fi

.PHONY: all clean build build-win64