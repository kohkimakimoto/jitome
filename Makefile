.PHONY: default build bindata deps

default: build

build: bindata
	go build .

bindata:
	@cd ./bindata/ && go-bindata -pkg bindata -prefix "src" -o bindata.go src/...

test:
	@go test -v -cover $$(go list ./... | grep -v vendor)

deps:
	glide up
