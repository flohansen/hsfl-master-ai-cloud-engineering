# Load Balancer

This is a simple load balancer that distributes requests to multiple servers.

It supports the following load balancing algorithms:

- Round Robin
- Random
- IP Hash

## Run

The load balancer only works within Docker, because it uses Docker networks to connect to the servers and IP addresses to identify the servers. You can use this docker compose example to start the load-balancer:

```docker
services:
  load-balancer:
    build:
      context: ./
      dockerfile: ./src/load-balancer/Dockerfile
    command:
      [
        "/app/src/load-balancer/main",
        "-config",
        "/app/config/load-balancer.json"
      ]
    ports:
      - '3001:3000'
    environment:
      - HTTP_SERVER_PORT=3000
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - bulletinboard
networks:
  bulletinboard:
    driver: bridge
    name: bulletinboard
```

When running locally, make sure you can access the containers from the host.

## Config

The Load Balancer can be configured by environment variables and a `config.json` file.

The following environment variables are available:

| Name | Description | Default |
| ---- | ----------- | -------- |
| HTTP_SERVER_PORT | Port the http server listens on | 3000 |

The default `config.json` looks like this:

```json
{
  "image": "docker.io/yeasy/simple-web",
  "networkName": "bridge",
  "replicas": 3,
  "healthCheckIntervalSeconds": 5,
  "healthCheckPath": "/",
  "port": 80,
  "algorithm": "round_robin"
}
```

- Image: The image to use for the servers.
- NetworkName: The name of the network the servers are connected to.
- Replicas: The number of servers to start.
- HealthCheckIntervalSeconds: The interval in seconds in which the health of the servers is checked.
- HealthCheckPath: The path to use for the health check.
- Port: The port the servers are listening on.
- Algorithm: The load balancing algorithm to use. Supported values are `round_robin`, `random`, `ip_hash` and `least_connections`.

You may also use a `.env` file to set environment variables. See `.env.example` for an example.