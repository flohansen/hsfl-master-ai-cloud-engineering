# Stage 1: Build Image
FROM golang:1.21-alpine AS BuildStage

WORKDIR /app
COPY ./lib ./lib
COPY ./src/user-service ./src/user-service

WORKDIR /app/src/user-service
RUN go mod tidy
RUN go build -o /user-service

# Stage 2: Service Image
FROM scratch
WORKDIR /

COPY --from=BuildStage /user-service /user-service

EXPOSE 3003
CMD ["/user-service"]