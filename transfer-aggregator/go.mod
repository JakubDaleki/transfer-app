module github.com/JakubDaleki/transfer-app/transfer-aggregator

go 1.20

require (
	github.com/JakubDaleki/transfer-app/shared-dependencies v0.0.0
	github.com/segmentio/kafka-go v0.4.40
	google.golang.org/grpc v1.55.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/JakubDaleki/transfer-app/shared-dependencies v0.0.0 => ../shared-dependencies
