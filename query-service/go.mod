module github.com/JakubDaleki/transfer-app/query-service

go 1.19

require (
	github.com/JakubDaleki/transfer-app/shared-dependencies v0.0.0
	github.com/hashicorp/go-memdb v1.3.4
	google.golang.org/grpc v1.55.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/JakubDaleki/transfer-app/shared-dependencies v0.0.0 => ../shared-dependencies
