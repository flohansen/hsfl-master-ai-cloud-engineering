# Stage 1: Build Image
FROM golang:1.21-alpine AS BuildStage
RUN apk add --update nodejs npm

WORKDIR /app
COPY ./src/web-service ./src/web-service

WORKDIR /app/src/web-service/frontend
RUN npm ci
RUN npm run build

WORKDIR /app/src/web-service
RUN go build -o /web-service

# Stage 2: Service Image
FROM scratch
WORKDIR /

COPY --from=BuildStage /web-service /web-service
COPY --from=BuildStage /app/src/web-service/frontend/static /frontend/static

EXPOSE 3000
CMD ["/web-service"]