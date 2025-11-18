FROM golang:1.25.3-alpine AS builder
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/grpc_http_api ./cmd/main

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /out/grpc_http_api ./grpc_http_api
COPY ./api ./api
COPY ./db ./db

EXPOSE 8080
EXPOSE 5001

ENTRYPOINT ["./grpc_http_api"]
