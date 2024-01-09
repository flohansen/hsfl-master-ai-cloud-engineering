# Stage 1: Build Image
FROM golang:1.21-alpine AS BuildStage

WORKDIR /app
COPY ./lib ./lib
COPY ./src/product-service ./src/product-service

WORKDIR /app/src/product-service
RUN go mod tidy
RUN go build -o /product-service

# Stage 2: Service Image
FROM scratch
WORKDIR /

COPY --from=BuildStage /product-service /product-service

EXPOSE 3003
CMD ["/product-service"]