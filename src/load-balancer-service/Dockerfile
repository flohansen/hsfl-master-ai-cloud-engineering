FROM golang:1.21-alpine
ARG LOAD_BALANCER_CONFIG_PATH
ENV LOAD_BALANCER_CONFIG_PATH $LOAD_BALANCER_CONFIG_PATH

WORKDIR /app
COPY ./src/load-balancer-service ./src/load-balancer-service

WORKDIR /app/src/load-balancer-service
RUN go mod tidy
RUN go build -o ./main

EXPOSE 3000
CMD ["sh","-c","/app/src/load-balancer-service/main -config ${LOAD_BALANCER_CONFIG_PATH}"]