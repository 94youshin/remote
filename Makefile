fmt:
	@go fmt ./...
build:
	@CGO_ENABLE=1 go build -o remote cmd/remote.go
clean:
	@rm -f ./remote
.PHONY: build fmt clean