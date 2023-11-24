FROM golang:1.21-alpine

ARG PROXY_DOCKER_MODE
ENV PROXY_DOCKER_MODE $PROXY_DOCKER_MODE

WORKDIR /app
COPY ./lib ./lib
COPY ./src/http-proxy-service ./src/http-proxy-service

WORKDIR /app/src/http-proxy-service
RUN go mod tidy
RUN go build -o ./main

EXPOSE 8080
CMD ["/app/src/http-proxy-service/main"]