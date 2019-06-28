build:
	@go build -mod=vendor -o ./dist/sync ./cmd/sync

build-linux:
	@GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./dist/sync-linux-amd64 ./cmd/sync

test:
	@go test -mod=vendor -cover ./...

clean:
	@rm -rf ./dist

.PHONY: build build-linux test clean
