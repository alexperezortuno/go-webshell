BINARY=go_webshell
VERSION=1.0.0
BUILD_DIR=./build
BUILD_TIME=`date +%FT%T%z`
GOX_OS_ARCH="darwin/amd64 darwin/arm64 linux/386 linux/amd64 windows/386 windows/amd64"

.PHONY: default
default: build

.PHONY: directory
directory:
	mkdir -p ${BUILD_DIR}

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: run
run:
	export GOPATH=${GOPATH_1_20_2} && \
	go run cmd/api/main.go

.PHONY: build
build:
	export GOPATH=${GOPATH_1_20_2} &&
	mkdir -p ${BUILD_DIR} && \
	cp -r ./static ${BUILD_DIR} && \
	go build -a -o ${BUILD_DIR}/${BINARY} cmd/api/main.go

.PHONY: build-version
build-version:
	export GOPATH=${GOPATH_1_20_2} && \
	mkdir -p ${BUILD_DIR} && \
	go build -a -o ${BUILD_DIR}/${BINARY}-${VERSION} cmd/api/main.go

.PHONY: build-linux
build-linux:
	export GOPATH=${GOPATH_1_20_2} && \
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -ldflags "-X main.Version=${VERSION}" -a -o ${BUILD_DIR}/${BINARY}-${VERSION} cmd/api/main.go

.PHONY: deps
deps:
	dep ensure;

.PHONY: test
test:
	go test
