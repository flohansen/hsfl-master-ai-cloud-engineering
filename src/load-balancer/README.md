# Load Balancer

This is a simple load balancer that distributes requests to multiple servers.

It supports the following load balancing algorithms:

- Round Robin
- Random
- IP Hash

## Config

Config values are set via environment variables. The following environment variables are available:

| Name | Description |
| ---- | ----------- |
| HTTP_SERVER_PORT | Port the http server listens on |
| IMAGE | Image to use for the load balancer |
| NETWORK_NAME | Name of the network the containers are connected to |
| REPLICAS | Number of replicas to start |
| HEALTH_CHECK_INTERVAL_SECONDS | Interval in seconds to check the health of the containers |



You may also use a `.env` file to set environment variables. See `.env.example` for an example.