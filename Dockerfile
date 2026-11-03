FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go mod download \
    && main_package="$(go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./... | grep -v '^$' | head -n 1)" \
    && test -n "$main_package" \
    && go build -trimpath -ldflags="-s -w" -o /app "$main_package"

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app /app/server

EXPOSE 8080

ENTRYPOINT ["/app/server"]
