FROM golang:1.24.5-alpine3.22 AS builder_auth

COPY . /github.com/evgeniySeleznev/auth-project/source
WORKDIR /github.com/evgeniySeleznev/auth-project/source

RUN go mod download
RUN go build -o ./bin/auth_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder_auth /github.com/evgeniySeleznev/auth-project/source/bin/auth_server .

CMD ["./auth_server"]