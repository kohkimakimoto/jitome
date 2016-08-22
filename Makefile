.PHONY: default build bindata deps_update

default: build

build: bindata
	go build .

bindata:
	@cd ./bindata/ && go-bindata -pkg bindata -prefix "src" -o bindata.go src/...

deps_update:
	rm -rf vendor
	gom install
	rm -rf vendor/**/**/.git
	rm -rf vendor/**/**/**/.git
