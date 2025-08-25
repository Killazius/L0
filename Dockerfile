FROM golang:1.25.0-alpine AS builder

WORKDIR /cmd
COPY go.mod go.sum ./
RUN go mod download
COPY . .


FROM builder AS app-builder
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app ./cmd/app/

FROM builder AS migrator-builder
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o migrator ./cmd/migrator/



FROM alpine:3.22 AS app
WORKDIR /app
COPY --from=app-builder /cmd/app .
RUN mkdir -p /config
COPY --from=app-builder /cmd/config/ ./config/
COPY --from=app-builder /cmd/.env .
CMD ["./app"]


FROM alpine:3.22 AS migrator
WORKDIR /app
COPY --from=migrator-builder /cmd/migrator .
RUN mkdir -p /config
COPY --from=migrator-builder /cmd/config/ ./config/
COPY --from=migrator-builder /cmd/.env .
COPY --from=migrator-builder /cmd/migrations/ ./migrations/
CMD ["./migrator", "-command", "up"]