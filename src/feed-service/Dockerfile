FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/feed-service ./src/feed-service

WORKDIR /app/src/feed-service
RUN go build -o ./main

CMD ["/app/src/feed-service/main"]