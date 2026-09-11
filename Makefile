.PHONY: all run fmt fmt-html fmt-go test

all: run

run:
	go run .

fmt-html:
	npx prettier --write "templates/**/*.html"

fmt-go:
	go fmt ./...

fmt: fmt-html fmt-go

test:
	go test ./...
