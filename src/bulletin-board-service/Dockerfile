FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/bulletin-board-service ./src/bulletin-board-service

WORKDIR /app/src/bulletin-board-service
RUN go build -o ./main

CMD ["/app/src/bulletin-board-service/main"]