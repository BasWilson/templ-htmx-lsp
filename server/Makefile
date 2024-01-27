clean:
	go clean
	rm -r -f ./bin

compile: 
	@make clean
	mkdir ./bin
	go build -o ./bin ./cmd/...

compile_linux: 
	@make clean
	mkdir ./bin
	GOARCH=amd64 GOOS=linux go build -o ./bin ./cmd/...

run: 
	@make clean
	@make compile
	./bin/lsp

dev:
	@go run cmd/lsp/main.go