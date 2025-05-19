
.PHONY: build test fmt clean

build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin ./...

test:
	go test ./...

fmt:
	go fmt ./... && go vet ./... && golangci-lint run

clean:
	rm -rf ./bin/*