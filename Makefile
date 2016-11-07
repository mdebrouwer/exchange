PACKAGE := github.com/mdebrouwer/exchange
MAPPED_PATH := /workspace/src/$(PACKAGE)
GOLANG := docker run --rm -i -v "`pwd`:$(MAPPED_PATH)" -e "GOPATH=/workspace" -w "$(MAPPED_PATH)" golang:1.7.3
SOURCE := $(shell find . -name '*.go' | grep -v vendor)

.PHONY: clean test

default: vet test build

clean:
	rm -rf ./tools ./.deps ./vendor ./service/bindata.go

vet:
	$(GOLANG) go tool vet $(SOURCE)

test: tools/ginkgo
	$(GOLANG) ./tools/ginkgo -r -race

build: exchange

exchange: $(SOURCE) service/bindata.go
	$(GOLANG) go build

service/bindata.go: tools/go-bindata static/index.html
	./tools/go-bindata -pkg service -o service/bindata.go static/

tools/courier:
	mkdir -p tools
	curl --fail -L https://github.com/optiver/courier/releases/download/1.1.0/courier.gz | gunzip - > ./tools/courier
	chmod +x ./tools/courier

tools/ginkgo: .deps
	mkdir -p tools
	$(GOLANG) go build -o ./tools/ginkgo ./vendor/github.com/onsi/ginkgo/ginkgo

tools/go-bindata: .deps
	mkdir -p tools
	$(GOLANG) go build -o ./tools/go-bindata ./vendor/github.com/jteeuwen/go-bindata/go-bindata

.deps: tools/courier pins.json
	rm -rf vendor
	./tools/courier -reproduce
	touch .deps
