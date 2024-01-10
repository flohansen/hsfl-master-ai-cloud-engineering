# Data Service

The data services provides a REST API to create, read, update and delete bulletin board data.

## Config

Config values are set via environment variables. The following environment variables are available:

| Name | Description |
| ---- | ----------- |
| HTTP_SERVER_PORT | Port the http server listens on |
| AUTH_SERVICE_URL_GRPC | URL of the authentication service gRPC endpoint |
| DB_HOST | Hostname of PostgreSQL database |
| DB_PORT | Port of PostgreSQL database |
| DB_USER | Username of PostgreSQL database |
| DB_PASSWORD | Password of PostgreSQL database |
| DB_NAME | Name of PostgreSQL database |

You may also use a `.env` file to set environment variables. See `.env.example` for an example.