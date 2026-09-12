build:
	@go run main/main.go

dev:
	@go run main/main.go --watch

test:
	@go test ./...

vet:
	@go vet ./...

clean:
	@go clean -testcache
