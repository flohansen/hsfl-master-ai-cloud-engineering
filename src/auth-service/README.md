# Auth Service

The auth service is responsible for managing users. It provides a REST API to create, read, update and delete users.

## Database

The auth service uses a PostgreSQL database to store user data.

Setup postgres database using docker:

```bash
docker run --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
```

### Generate key

The auth service uses a ECDSA private key to sign JWTs. The key can be generated using the following command:

```bash
ssh-keygen -t ecdsa -f /path/to/key -m pem
```

## Config

Config values are set via environment variables. The following environment variables are available:

| Name | Description |
| ---- | ----------- |
| PORT | Port the auth service listens on |
| JWT_SIGN_KEY | Path to ECDSA private key used to sign JWTs |
| DB_HOST | Hostname of PostgreSQL database |
| DB_PORT | Port of PostgreSQL database |
| DB_USER | Username of PostgreSQL database |
| DB_PASSWORD | Password of PostgreSQL database |
| DB_NAME | Name of PostgreSQL database |
