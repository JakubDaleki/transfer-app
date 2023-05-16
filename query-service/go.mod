module github.com/JakubDaleki/transfer-app/webapp

go 1.19


replace github.com/JakubDaleki/transfer-app/shared-dependencies => ../shared-dependencies
require (
	github.com/JakubDaleki/transfer-app/shared-dependencies v0.0.0-20230516102607-152713c33501
	github.com/hashicorp/go-memdb v1.3.4
)

require (
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
)
