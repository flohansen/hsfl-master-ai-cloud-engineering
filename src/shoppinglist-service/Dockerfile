FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/shoppinglist-service ./src/shoppinglist-service

WORKDIR /app/src/shoppinglist-service
RUN go mod tidy
RUN go build -o ./main

EXPOSE 3002
CMD ["/app/src/shoppinglist-service/main"]