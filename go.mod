module github.com/duckpie/bfb-security-microservice

go 1.18

require (
	github.com/duckpie/cherry v0.0.0-20220411202740-c286ab674c25
	github.com/spf13/cobra v1.4.0
	github.com/wrs-news/golang-proto v0.3.5
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/myesui/uuid v1.0.0 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	gopkg.in/stretchr/testify.v1 v1.2.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/oklog/oklog v0.3.2
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.1
	github.com/twinj/uuid v1.0.0
	google.golang.org/grpc v1.45.0
)

replace google.golang.org/genproto => ./libs/genproto
