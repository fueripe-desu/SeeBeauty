build:
	@go build -o bin/bkalpha main.go 

run: build
	@./bin/bkalpha

test:
	@go test -v ./...
