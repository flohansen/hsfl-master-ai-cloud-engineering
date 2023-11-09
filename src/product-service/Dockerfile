FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/product-service ./src/product-service

WORKDIR /app/src/product-service
RUN go mod tidy
RUN go build -o ./main

EXPOSE 3003
CMD ["/app/src/product-service/main"]