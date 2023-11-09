FROM golang:1.21-alpine
RUN apk add --update nodejs npm

WORKDIR /app
COPY ./src/web-service ./src/web-service

WORKDIR /app/src/web-service/frontend
RUN npm ci
RUN npm run build

WORKDIR /app/src/web-service
RUN go build -o ./main

EXPOSE 3000
CMD ["/app/src/web-service/main"]