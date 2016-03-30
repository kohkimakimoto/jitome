.PHONY: default build bindata deps updatedeps

default: build

build:
	go build .

bindata:
	@cd ./bindata/src/ && go-bindata -pkg bindata -ignore "^\." -o ../bindata.go ./...

deps:
	gom install

updatedeps:
	rm Gomfile.lock; rm -rf vendor; gom install && gom lock
