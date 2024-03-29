BUILD_ENV := CGO_ENABLED=0
ESPACK_VERSION = $(shell cat version.txt)
BUILD = `date +%FT%T%z`
LDFLAGS = -ldflags "-w -s -X main.Version=${ESPACK_VERSION} -X main.Build=${BUILD_ENV}"

TARGET_EXEC := espack

.PHONY:clean setup build-linux build-osx build-windows build-osx-arm64

all:clean setup build-linux build-osx build-windows build-osx-arm64

run:

check:
	go fmt ./
	go vet ./

clean:
	go clean
	rm -rf build

setup:
	mkdir -p build/linux
	mkdir -p build/osx/amd64
	mkdir -p build/osx/arm64
	mkdir -p build/windows

build-linux: setup
	GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/linux/${TARGET_EXEC} ./cmd

build-osx: setup
	GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/osx/amd64/${TARGET_EXEC} ./cmd

build-osx-arm64: setup
	GOARCH=arm64 GOOS=darwin go build ${LDFLAGS} -o build/osx/arm64/${TARGET_EXEC} ./cmd

build-windows: setup
	GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/windows/${TARGET_EXEC}.exe ./cmd

help:
	@echo "make 格式化go代码 并编译生成二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make check 格式化go代码"
	@echo "make run 直接运行程序"