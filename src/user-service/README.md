# The user service component
[![Run tests (user-service)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-user-service.yml/badge.svg)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-user-service.yml)
The user service component provides user data and helper functions.

## Configuration
You can configure this service with environmental variables or files (.env)

### Example configuration
```dotenv
RQLITE_HOST="dbhost"
RQLITE_PORT=4001
RQLITE_USER="<user-name>"
RQLITE_PASSWORD="<password>"

HTTP_SERVER_PORT=3001

GRPC_SERVER_PORT=50051
```