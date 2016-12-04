PACKAGE := github.com/mdebrouwer/exchange
MAPPED_PATH := /workspace/src/$(PACKAGE)

DOCKER := docker run --rm -i -v "`pwd`:$(MAPPED_PATH)" -e "GOPATH=/workspace" -w "$(MAPPED_PATH)"
GOLANG := $(DOCKER) golang:1.7.3
NODE := $(DOCKER) node:7.1.0
BASH := $(DOCKER) centos:7

GO_SOURCE := $(shell $(BASH) bash -c "find . -name '*.go' | grep -v vendor")
JS_SOURCE := $(shell $(BASH) find static -name '*.js')
LESS_SOURCE := $(shell $(BASH) find static -name '*.less')

.PHONY: clean vet test build keys

default: vet test build

clean:
	$(BASH) rm -rf ./tools ./.go_deps ./.js_deps ./vendor ./node_modules ./bundle ./service/bindata.go ./keygen/keygen

dev: .js_deps exchange
	docker-compose -f dev.yml up

vet: .go_deps
	$(GOLANG) go tool vet $(GO_SOURCE)

test: tools/ginkgo .go_deps
	$(GOLANG) ./tools/ginkgo -r -race

build: exchange

exchange: $(GO_SOURCE) service/bindata.go .go_deps
	$(GOLANG) go build

service/bindata.go: tools/go-bindata bundle/index.html bundle/assets/bundle.js
	$(BASH) ./tools/go-bindata -pkg service -o service/bindata.go bundle/...

bundle/index.html: static/index.html
	$(BASH) mkdir -p bundle
	$(BASH) cp static/index.html bundle/index.html

bundle/assets/bundle.js: $(JS_SOURCE) $(LESS_SOURCE) .js_deps webpack.config.js
	$(BASH) mkdir -p bundle/assets
	$(NODE) ./node_modules/.bin/webpack

tools/courier:
	$(BASH) mkdir -p tools
	$(DOCKER) tutum/curl bash -c "curl --fail -L https://github.com/optiver/courier/releases/download/1.1.0/courier.gz | gunzip - > ./tools/courier"
	$(BASH) chmod +x ./tools/courier

tools/ginkgo: .go_deps
	$(BASH) mkdir -p tools
	$(GOLANG) go build -o ./tools/ginkgo ./vendor/github.com/onsi/ginkgo/ginkgo

tools/go-bindata: .go_deps
	$(BASH) mkdir -p tools
	$(GOLANG) go build -o ./tools/go-bindata ./vendor/github.com/jteeuwen/go-bindata/go-bindata

.js_deps: package.json
	$(BASH) rm -rf node_modules
	$(NODE) npm install
	$(BASH) touch .js_deps

.go_deps: tools/courier pins.json
	$(BASH) rm -rf vendor
	$(GOLANG) ./tools/courier -reproduce
	$(BASH) touch .go_deps

keys: keygen/keygen
	$(BASH) ./keygen/keygen

keygen/keygen: .go_deps keygen/main.go
	$(GOLANG) go build  -o ./keygen/keygen ./keygen/main.go
