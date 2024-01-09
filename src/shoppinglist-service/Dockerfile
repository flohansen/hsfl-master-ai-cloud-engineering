# Stage 1: Build Image
FROM golang:1.21-alpine AS BuildStage

WORKDIR /app
COPY ./lib ./lib
COPY ./src/shoppinglist-service ./src/shoppinglist-service

WORKDIR /app/src/shoppinglist-service
RUN go mod tidy
RUN go build -o /shoppinglist-service

# Stage 2: Service Image
FROM scratch
WORKDIR /

COPY --from=BuildStage /shoppinglist-service /shoppinglist-service

EXPOSE 3002
CMD ["/shoppinglist-service"]