module hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service

go 1.21

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	golang.org/x/crypto v0.14.0 // indirect
)
replace (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/router" => ./../../lib/router
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc" => ./../../lib/rpc
)