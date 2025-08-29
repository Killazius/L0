.PHONY: produce docker test lint swag
COUNT ?= 1

produce:
	go run cmd/kafka/producer.go -m $(COUNT)

docker:
	docker compose down && docker image prune -f && docker compose up -d --build

test:
	go test -v -race -parallel 5 -shuffle=on -coverprofile=./cover.out -covermode=atomic ./...

lint:
	golangci-lint run ./...

swag:
	swag init -g ./cmd/app/main.go -o ./docs