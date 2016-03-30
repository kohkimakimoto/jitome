.PHONY: default build bindata deps updatedeps

default: build

build:
	go build .

bindata:
	@cd ./bindata/src/ && go-bindata -pkg bindata -ignore "^\." -o ../bindata.go ./...

deps:
	rm -rf vendor
	gom install
	rm -rf vendor/**/**/.git
	rm -rf vendor/**/**/**/.git