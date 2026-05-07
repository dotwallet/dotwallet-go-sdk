module github.com/dotwallet/dotwallet-go-sdk

go 1.25.0

require (
	github.com/btcsuite/btcd v0.24.0
	github.com/btcsuite/btcutil v1.0.2
	github.com/go-resty/resty/v2 v2.17.2
	github.com/jarcoal/httpmock v1.4.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mrz1836/go-api-router v1.0.14
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/matryer/respond v1.0.1 // indirect
	github.com/mrz1836/go-logger v1.0.6 // indirect
	github.com/mrz1836/go-parameters v1.0.8 // indirect
	github.com/newrelic/go-agent/v3 v3.43.3 // indirect
	github.com/newrelic/go-agent/v3/integrations/nrhttprouter v1.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260504160031-60b97b32f348 // indirect
	google.golang.org/grpc v1.81.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Issue with using wrong version of Redigo
replace github.com/btcsuite/btcd => github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3
