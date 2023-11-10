FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/web-service ./src/web-service

WORKDIR /app/src/web-service
RUN go build -o ./main

CMD ["/app/src/web-service/main"]