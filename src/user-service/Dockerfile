FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/user-service ./src/user-service

WORKDIR /app/src/user-service
RUN go mod tidy
RUN go build -o ./main

EXPOSE 3003
CMD ["/app/src/user-service/main"]