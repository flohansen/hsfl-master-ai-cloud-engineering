module github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service

go 1.21.3

require (
	github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib v0.0.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/stretchr/testify v1.8.4
	golang.org/x/crypto v0.14.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

replace github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib => ../../lib
