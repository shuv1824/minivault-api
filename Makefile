run: build
	@./bin/minivault

build:
	@go build -o bin/minivault .
