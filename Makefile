docker:
	docker compose down && docker image prune -f && docker compose up -d --build

test:
	go test -v -race -parallel 5 -shuffle=on -coverprofile=./cover.out -covermode=atomic ./...

lint:
	golangci-lint run ./...
