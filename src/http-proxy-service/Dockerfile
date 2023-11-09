FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/http-proxy-service ./src/http-proxy-service

WORKDIR /app/src/http-proxy-service
RUN go mod tidy
RUN go build -o ./main

EXPOSE 8080
CMD ["/app/src/http-proxy-service/main"]