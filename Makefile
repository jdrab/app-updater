
LINUX_BIN="app-updater"
LINUX_SERVICE="my-service"
LINUX_APP="my-app"

WIN_BIN="app-updater.exe"
WIN_SERVICE="my-service"
WIN_APP="Application.exe"

#main output directory
OUT=out

VERSION=`git describe --tags --long`

MAKEFLAGS += --silent

# Builds the project
all:	build-linux build-win64 build-win32

build-linux: 
	export GOOG=linux
	export GOARCH=amd64
	echo "building for linux/amd64"
	CGO_ENABLED=0 go build -a -o ${OUT}/linux/amd64/${LINUX_BIN} -ldflags "-w -s -X main.Version=${VERSION} -X main.runtimeApp=${LINUX_SERVICE} -X main.runtimeService=${LINUX_APP}"

build-win64:
	export GOOG=windows
	export GOARCH=amd64
	echo "building for windows/amd64"
	CGO_ENABLED=0 go build -a -o ${OUT}/windows/amd64/${WIN_BIN} -ldflags "-w -s -X main.Version=${VERSION} -X main.runtimeApp=${WIN_SERVICE} -X main.runtimeService=${WIN_APP}"

build-win32:
	export GOOG=windows
	export GOARCH=386
	echo "building for windows/386"
	CGO_ENABLED=0 go build -a -o ${OUT}/windows/386/${WIN_BIN} -ldflags "-w -s -X main.Version=${VERSION} -X main.runtimeApp=${WIN_SERVICE} -X main.runtimeService=${WIN_APP}"



# Cleans our project: deletes binaries
clean:
	if [ -d ${OUT} ];then echo "removing directory ${OUT}"; rm -rf ${OUT} ;fi

.PHONY: clean build-linux