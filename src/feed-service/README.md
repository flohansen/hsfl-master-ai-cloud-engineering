# Feed Service

The feed service is responsible for managing feeds. It provides a REST API to fetch different types of feeds. For example a personal feed based on certain topics or a chronological feed.

## Config

Config values are set via environment variables. The following environment variables are available:

| Name | Description |
| ---- | ----------- |
| HTTP_SERVER_PORT | Port the http server listens on |

You may also use a `.env` file to set environment variables. See `.env.example` for an example.