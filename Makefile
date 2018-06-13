.PHONY: all
all: dep build

.PHONY: dep
dep:
	docker-compose run go dep ensure -vendor-only

.PHONY: build
build:
	go build -o ./bin/go-sqsd main.go

.PHONY: build-linux
build-linux:
	docker-compose run go go build -o ./bin/go-sqsd main.go
