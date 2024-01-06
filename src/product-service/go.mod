module hsfl.de/group6/hsfl-master-ai-cloud-engineering/product-service

go 1.21

require (
	github.com/caarlos0/env/v10 v10.0.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.8.4
	golang.org/x/sync v0.4.0
	google.golang.org/grpc v1.60.0
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/config v0.0.0-00010101000000-000000000000
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router v0.0.0-00010101000000-000000000000
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/config => ./../../lib/config
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router => ./../../lib/router
	hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc => ./../../lib/rpc
)
