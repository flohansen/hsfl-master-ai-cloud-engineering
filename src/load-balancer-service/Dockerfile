FROM golang:1.21-alpine
ARG IMAGE
ARG REPLICAS
ARG NETWORK_NAME
ENV IMAGE $IMAGE
ENV REPLICAS $REPLICAS
ENV NETWORK_NAME $NETWORK_NAME

WORKDIR /app
COPY ./src/load-balancer-service ./src/load-balancer-service

WORKDIR /app/src/load-balancer-service
RUN go mod tidy
RUN go build -o ./main

EXPOSE 3000
CMD ["sh","-c","/app/src/load-balancer-service/main -image ${IMAGE} -replicas ${REPLICAS} -network ${NETWORK_NAME}"]