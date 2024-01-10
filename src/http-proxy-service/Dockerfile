# Stage 1: Build Image
FROM golang:1.21-alpine AS BuildStage

WORKDIR /app
COPY ./lib ./lib
COPY ./src/http-proxy-service ./src/http-proxy-service

WORKDIR /app/src/http-proxy-service
RUN go mod tidy
RUN go build -o /http-proxy-service

# Stage 2: Service Image
FROM scratch
WORKDIR /

COPY --from=BuildStage /http-proxy-service /http-proxy-service

EXPOSE 8080
CMD ["/http-proxy-service"]