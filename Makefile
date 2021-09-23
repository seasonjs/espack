BUILD_ENV := CGO_ENABLED=0
ESBUILD_VERSION = $(shell cat version.txt)
BUILD = `date +%FT%T%z`
LDFLAGS = -ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

TARGET_EXEC := espack

.PHONY:all clean setup build-linux build-osx build-windows

all:clean setup build-linux build-osx build-windows

clean:
	rm -rf build

setup:
	mkdir -p build/linux
	mkdir -p build/osx
	mkdir -p build/windows

build-linux: setup
	GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/linux/${TARGET_EXEC}

build-osx: setup
	GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/osx/${TARGET_EXEC}

build-windows: setup
	GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/windows/${TARGET_EXEC}.exe