PACKAGE := github.com/mdebrouwer/exchange
MAPPED_PATH := /workspace/src/$(PACKAGE)
GOLANG := docker run --rm -it -v "`pwd`:$(MAPPED_PATH)" -e "GOPATH=/workspace" -w "$(MAPPED_PATH)" golang:1.7.3
GO_SOURCE := $(shell find . -name '*.go' | grep -v vendor)
NODE := docker run --rm -it -v "`pwd`:/workspace" -w "/workspace" node:7.1.0
JS_SOURCE := $(shell find static -name '*.js')
LESS_SOURCE := $(shell find static -name '*.less')

.PHONY: clean vet test build keys

default: vet test build

clean:
	rm -rf ./tools ./.go_deps ./.js_deps ./vendor ./node_modules ./bundle ./service/bindata.go

dev: .js_deps
	# Not using docker container: https://github.com/nodejs/node/issues/4182
	./node_modules/.bin/webpack-dev-server --content-base static/ --host 0.0.0.0 --watch-poll

vet: .go_deps
	$(GOLANG) go tool vet $(GO_SOURCE)

test: tools/ginkgo .go_deps
	$(GOLANG) ./tools/ginkgo -r -race

build: exchange

exchange: $(GO_SOURCE) service/bindata.go .go_deps
	$(GOLANG) go build

service/bindata.go: tools/go-bindata bundle/index.html bundle/assets/bundle.js
	./tools/go-bindata -pkg service -o service/bindata.go bundle/...

bundle/index.html: static/index.html
	mkdir -p bundle
	cp static/index.html bundle/index.html

bundle/assets/bundle.js: $(JS_SOURCE) $(LESS_SOURCE) .js_deps webpack.config.js
	mkdir -p bundle/assets
	$(NODE) ./node_modules/.bin/webpack

tools/courier:
	mkdir -p tools
	curl --fail -L https://github.com/optiver/courier/releases/download/1.1.0/courier.gz | gunzip - > ./tools/courier
	chmod +x ./tools/courier

tools/ginkgo: .go_deps
	mkdir -p tools
	$(GOLANG) go build -o ./tools/ginkgo ./vendor/github.com/onsi/ginkgo/ginkgo

tools/go-bindata: .go_deps
	mkdir -p tools
	$(GOLANG) go build -o ./tools/go-bindata ./vendor/github.com/jteeuwen/go-bindata/go-bindata

.js_deps: package.json
	rm -rf node_modules
	$(NODE) npm install
	touch .js_deps

.go_deps: tools/courier pins.json
	rm -rf vendor
	./tools/courier -reproduce
	touch .go_deps

keys: keygen/keygen
	./keygen/keygen

keygen/keygen:
	$(GOLANG) go build  -o ./keygen/keygen ./keygen/main.go