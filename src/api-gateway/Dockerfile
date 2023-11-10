FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/api-gateway ./src/api-gateway

WORKDIR /app/src/api-gateway
RUN go build -o ./main

CMD ["/app/src/api-gateway/main"]