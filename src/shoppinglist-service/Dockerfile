FROM golang:1.21-alpine
RUN apk update && apk add --no-cache make protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app
COPY ./lib ./lib
COPY ./src/shoppinglist-service ./src/shoppinglist-service

WORKDIR /app/src/shoppinglist-service
RUN go mod tidy
RUN go generate
RUN go build -o ./main

EXPOSE 3002
CMD ["/app/src/shoppinglist-service/main"]