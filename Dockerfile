FROM golang:1.25.3 AS builder
WORKDIR /service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o ./bin/service ./cmd/service

FROM alpine:latest
COPY --from=builder /service/bin/service /service
COPY --from=builder /service/config.env /config.env

EXPOSE 8080
ENTRYPOINT ["/service"]
