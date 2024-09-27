build:
	@go build -o bin/bkalpha  

run: build
	@./bin/bkalpha

test:
	@go test -v ./...
