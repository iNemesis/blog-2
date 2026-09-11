.PHONY: all run fmt fmt-html fmt-go test optimize optimize-lossy

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

optimize:
	go run tools/optimize.go

optimize-lossy:
	go run tools/optimize.go -lossy
