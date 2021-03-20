FROM golang:1.16.2-alpine3.13 AS builder
WORKDIR /build

RUN apk add protoc
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

COPY protofiles protofiles
RUN mkdir pb && \
    protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --proto_path=protofiles protofiles/*.proto

COPY go.mod go.sum ./
RUN go mod download

ENV CGO_ENABLED 0

COPY main.go .
RUN go build -o kv-client .


FROM alpine:3.13
COPY --from=builder /build/kv-client /usr/local/bin

ENTRYPOINT ["kv-client"]
