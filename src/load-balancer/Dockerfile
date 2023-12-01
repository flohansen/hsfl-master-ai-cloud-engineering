FROM golang:1.21-alpine

WORKDIR /app
COPY ./lib ./lib
COPY ./src/load-balancer ./src/load-balancer

WORKDIR /app/src/load-balancer
RUN go build -o ./main

CMD ["/app/src/load-balancer/main"]