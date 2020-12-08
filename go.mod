module secProxy

go 1.14

require (
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/gomodule/redigo v1.8.3
	github.com/google/uuid v1.1.2 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	google.golang.org/grpc v1.33.2 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
