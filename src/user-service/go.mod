module hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service

go 1.21

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/stretchr/testify v1.8.4
	golang.org/x/crypto v0.14.0
	google.golang.org/grpc v1.60.1
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router v0.0.0-00010101000000-000000000000
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router => ./../../lib/router
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc => ./../../lib/rpc
)
