build:
	@go build -o bin/demo main.go

run: build
	@./bin/demo