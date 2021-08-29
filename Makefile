TEST?=./...
HOSTNAME=registry.terraform.io
NAMESPACE=jaysonsantos
NAME=jumphost
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
CURRENT_OS:=$(shell uname | tr [A-Z] [a-z])
CURRENT_ARCH:=$(shell uname -m | sed 's/x86_64/amd64/')
OS_ARCH=$(CURRENT_OS)_$(CURRENT_ARCH)
GO_MOD_FILES := go.mod go.sum

default: install

.cache/go-dependencies: $(GO_MOD_FILES)
	go get -v ./...
	mkdir -p .cache
	touch .cache/go-dependencies

build: .cache/go-dependencies
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4 -v

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -v
