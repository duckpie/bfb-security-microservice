module github.com/duckpie/bfb-security-microservice

go 1.18

require (
	github.com/spf13/cobra v1.4.0
	github.com/wrs-news/golang-proto v0.3.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/protobuf v1.26.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/oklog/run v1.1.0 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/oklog/oklog v0.3.2
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/twinj/uuid v1.0.0
	google.golang.org/grpc v1.45.0
)

replace google.golang.org/genproto => ./libs/genproto
