PACKAGE := github.com/mdebrouwer/exchange
MAPPED_PATH := /workspace/src/$(PACKAGE)
GOLANG := docker run --rm -i -v "`pwd`:$(MAPPED_PATH)" -e "GOPATH=/workspace" -w "$(MAPPED_PATH)" golang:1.7.3
SOURCE := $(shell find . -name '*.go' | grep -v vendor)

.PHONY: clean test

default: vet test build

clean:
	rm -rf tools .deps vendor

vet:
	$(GOLANG) go tool vet $(SOURCE)

test: tools/ginkgo
	$(GOLANG) ./tools/ginkgo -r -race

build: exchange

exchange: $(SOURCE)
	$(GOLANG) go build

tools/courier:
	mkdir -p tools
	curl --fail -L https://github.com/optiver/courier/releases/download/1.1.0/courier.gz | gunzip - > ./tools/courier
	chmod +x ./tools/courier

tools/ginkgo: .deps
	mkdir -p tools
	$(GOLANG) go build -o ./tools/ginkgo ./vendor/github.com/onsi/ginkgo/ginkgo

.deps: tools/courier pins.json
	rm -rf vendor
	./tools/courier -reproduce
	touch .deps
