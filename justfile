set dotenv-load := true

app := "mabata"
main := "./cmd/mabata"

fmt:
	gofmt -w cmd internal

tidy:
	go mod tidy

build:
	mkdir -p bin
	go build -o bin/{{app}} {{main}}

run:
	go run {{main}}

test:
	go test ./...

clean:
	rm -rf bin

check: fmt test

dev: tidy run
