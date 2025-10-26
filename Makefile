.PHONY: run test build docker-up docker-down fmt

run:
	go run ./cmd/server

test:
	go test ./...

build:
	go build -o bin/server ./cmd/server

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

docker-up:
	docker compose --profile dev up --build

docker-down:
	docker compose down
