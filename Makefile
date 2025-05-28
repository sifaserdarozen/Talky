
.PHONY: build test fmt clean

VERSION:=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE:=$(shell date +%Y-%m-%dT%H:%M:%S)
DIRTY:=$(shell ! [ -z "`git status --porcelain=v1 2>/dev/null`" ] && echo "dirty-")
BUILD_VERSION:=${DIRTY}${VERSION}

LDFLAGS += -X "github.com/sifaserdarozen/Talky/talky.Version=$(BUILD_VERSION)"
LDFLAGS += -X "github.com/sifaserdarozen/Talky/talky.BuildDate=$(BUILD_DATE)"

build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin ./...

test:
	go test ./...

docker-build-asterisk:
	docker build -t local-asterisk -f docker/asterisk/Dockerfile .

docker-build: docker-build-asterisk

docker-run-asterisk:
	docker run \
	-it \
	-v ./docker/asterisk/conf:/etc/asterisk \
	-p 6060:5060/udp \
	-p 6060:5060/tcp \
	local-asterisk:latest

docker-compose-up:
	docker compose -f docker/sip-trunk/docker-compose.yaml up

fmt:
	go fmt ./... && go vet ./... && golangci-lint run

clean:
	rm -rf ./bin/*