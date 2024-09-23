module github.com/dotwallet/dotwallet-go-sdk

go 1.19

require (
	github.com/btcsuite/btcd v0.24.0
	github.com/btcsuite/btcutil v1.0.2
	github.com/go-resty/resty/v2 v2.15.2
	github.com/jarcoal/httpmock v1.3.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mrz1836/go-api-router v0.7.3
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/matryer/respond v1.0.1 // indirect
	github.com/mrz1836/go-logger v0.3.3 // indirect
	github.com/mrz1836/go-parameters v0.4.1 // indirect
	github.com/newrelic/go-agent/v3 v3.30.0 // indirect
	github.com/newrelic/go-agent/v3/integrations/nrhttprouter v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240308144416-29370a3891b7 // indirect
	google.golang.org/grpc v1.62.1 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Issue with using wrong version of Redigo
replace github.com/btcsuite/btcd => github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3
