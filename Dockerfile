FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOTOOLCHAIN=auto
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /app/server ./cmd/main

FROM alpine:3.21

RUN apk add --no-cache ca-certificates && adduser -D -H appuser
WORKDIR /app
COPY --from=builder /app/server ./server

USER appuser
EXPOSE 8080
ENTRYPOINT ["./server"]
