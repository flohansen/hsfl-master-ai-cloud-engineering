# API Gateway

The API gateway is responsible for routing requests to the appropriate service.

## Config

Config values are set via environment variables. The following environment variables are available:

| Name | Description |
| ---- | ----------- |
| HTTP_SERVER_PORT | Port the http server listens on |
| AUTH_SERVICE_URL | URL of auth service |
| BULLETIN_BOARD_SERVICE_URL | URL of bulletin board service |
| FEED_SERVICE_URL | URL of feed service |
| WEB_SERVICE_URL | URL of web service |

You may also use a `.env` file to set environment variables. See `.env.example` for an example.
