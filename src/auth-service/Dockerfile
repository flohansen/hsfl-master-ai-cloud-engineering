FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/auth-service ./src/auth-service

WORKDIR /app/src/auth-service
RUN go build -o ./main

CMD ["/app/src/auth-service/main"]