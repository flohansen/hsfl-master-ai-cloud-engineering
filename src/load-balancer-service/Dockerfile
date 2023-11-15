FROM golang:1.21-alpine

WORKDIR /app
COPY . ./src/load-balancer-service

WORKDIR /app/src/load-balancer-service
RUN go mod tidy
RUN go build -o ./main

EXPOSE 3000
CMD ["/app/src/load-balancer-service/main"]