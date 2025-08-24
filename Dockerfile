FROM golang:1.25.0-alpine AS builder

WORKDIR /cmd
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app ./cmd/app/

FROM alpine:3.22

WORKDIR /cmd
COPY --from=builder /cmd/app .
RUN mkdir -p /config
COPY --from=builder /cmd/config/ ./config/
COPY --from=builder /cmd/.env .
CMD ["./app"]