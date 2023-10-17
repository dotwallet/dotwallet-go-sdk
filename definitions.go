package dotwallet

import (
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/go-resty/resty/v2"
)

const (
	// Package configuration defaults
	apiVersion                   string = "v1"                           // Version of the API
	defaultHost                         = "https://api.ddpurse.com"      // Default host for API endpoints
	defaultHTTPTimeout                  = 10 * time.Second               // Default timeout for all GET requests in seconds
	defaultRefreshTokenExpiresIn        = 7 * 24 * time.Hour             // Default is 7 days (from documentation)
	defaultRetryCount            int    = 2                              // Default retry count for HTTP requests
	defaultUserAgent                    = "dotwallet-go-sdk: " + version // Default user agent
	version                      string = "v0.2.0"                       // dotwallet-go-sdk version

	// Grants
	grantTypeAuthorizationCode = "authorization_code"
	grantTypeClientCredentials = "client_credentials"
	grantTypeRefreshToken      = "refresh_token"

	// Endpoints
	authorizeURI          = "/" + apiVersion + "/oauth2/authorize"
	getAccessTokenURI     = "/" + apiVersion + "/oauth2/get_access_token"
	getUserInfo           = "/" + apiVersion + "/user/get_user_info"
	getUserReceiveAddress = "/" + apiVersion + "/user/get_user_receive_address"

	// NFT Endpoints
	getNft               = "/phoenix/user_all/zy_get_nft"
	mintNft              = "/phoenix/user_all/zy_mint_nft"
	transferNftToAddress = "/phoenix/user_all/zy_transfer_nft_to_address"

	// Merkle proof endpoints
	getMerkleProof    = "/phoenix/public/get_merkle_proof"
	getRawTransaction = "/phoenix/public/get_rawtransaction"

	// Headers
	headerAuthorization = "Authorization"

	// ScopeUserInfo is for getting user info
	ScopeUserInfo = "user.info"

	// ScopeAutoPayBSV is for auto-pay with a BSV balance
	ScopeAutoPayBSV = "autopay.bsv"

	// ScopeAutoPayBTC is for auto-pay with a BTC balance
	ScopeAutoPayBTC = "autopay.btc"

	// ScopeAutoPayETH is for auto-pay with a ETH balance
	ScopeAutoPayETH = "autopay.eth"
)

// NFT constants
const (
	NftVoutLen            = 2188
	NftVinPushedDataCount = 11
)

// NftVoutScriptPrefix is the prefix
var NftVoutScriptPrefix = []byte{81, 1, 64, 1, 118, 1, 136, 1, 169, 1, 172, 97, 94, 121, 97, 0, 121, 1, 104, 127, 119, 0, 0, 82, 121, 81, 127, 117, 0, 127, 119, 0, 121, 1, 253, 135, 99, 97, 83, 121, 83, 127, 117, 81, 127, 119, 0, 121, 1, 0, 126, 129, 81, 122, 117, 97, 83, 122, 117, 82, 122, 82, 122, 83, 121, 83, 84, 121, 147, 127, 117, 83, 127, 119, 82, 122, 117, 81, 122, 103, 0, 121, 1, 254, 135, 99, 97, 83, 121, 85, 127, 117, 81, 127, 119, 0, 121, 1, 0, 126, 129, 81, 122, 117, 97, 83, 122, 117, 82, 122, 82, 122, 83, 121, 85, 84, 121, 147, 127, 117, 85, 127, 119, 82, 122, 117, 81, 122, 103, 0, 121, 1, 255, 135, 99, 97, 83, 121, 89, 127, 117, 81, 127, 119, 0, 121, 1, 0, 126, 129, 81, 122, 117, 97, 83, 122, 117, 82, 122, 82, 122, 83, 121, 89, 84, 121, 147, 127, 117, 89, 127, 119, 82, 122, 117, 81, 122, 103, 97, 83, 121, 81, 127, 117, 0, 127, 119, 0, 121, 1, 0, 126, 129, 81, 122, 117, 97, 83, 122, 117, 82, 122, 82, 122, 83, 121, 81, 84, 121, 147, 127, 117, 81, 127, 119, 82, 122, 117, 81, 122, 104, 104, 104, 81, 121, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 81, 122, 117, 97, 0, 121, 130, 119, 0, 121, 1, 20, 148, 82, 121, 81, 121, 127, 119, 83, 121, 82, 121, 127, 117, 0, 1, 18, 121, 0, 160, 99, 97, 1, 19, 121, 90, 121, 89, 121, 126, 1, 20, 126, 81, 121, 126, 90, 121, 126, 88, 121, 126, 81, 122, 117, 97, 97, 0, 121, 1, 20, 121, 0, 121, 88, 128, 97, 82, 121, 0, 121, 130, 119, 0, 81, 121, 2, 253, 0, 159, 99, 97, 81, 121, 81, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 81, 122, 117, 103, 81, 121, 3, 0, 0, 1, 159, 99, 1, 253, 97, 82, 121, 82, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 103, 81, 121, 5, 0, 0, 0, 0, 1, 159, 99, 1, 254, 97, 82, 121, 84, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 103, 81, 121, 9, 0, 0, 0, 0, 0, 0, 0, 0, 1, 159, 99, 1, 255, 97, 82, 121, 88, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 104, 104, 104, 104, 0, 121, 83, 121, 126, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 81, 122, 117, 97, 82, 121, 81, 121, 126, 83, 122, 117, 82, 122, 82, 122, 82, 121, 117, 117, 117, 104, 96, 121, 0, 160, 99, 81, 121, 1, 18, 121, 126, 97, 0, 121, 1, 18, 121, 0, 121, 88, 128, 97, 82, 121, 0, 121, 130, 119, 0, 81, 121, 2, 253, 0, 159, 99, 97, 81, 121, 81, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 81, 122, 117, 103, 81, 121, 3, 0, 0, 1, 159, 99, 1, 253, 97, 82, 121, 82, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 103, 81, 121, 5, 0, 0, 0, 0, 1, 159, 99, 1, 254, 97, 82, 121, 84, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 103, 81, 121, 9, 0, 0, 0, 0, 0, 0, 0, 0, 1, 159, 99, 1, 255, 97, 82, 121, 88, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 104, 104, 104, 104, 0, 121, 83, 121, 126, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 81, 122, 117, 97, 82, 121, 81, 121, 126, 83, 122, 117, 82, 122, 82, 122, 82, 121, 117, 117, 117, 104, 94, 121, 0, 160, 99, 97, 95, 121, 90, 121, 89, 121, 126, 1, 20, 126, 81, 121, 126, 90, 121, 126, 88, 121, 126, 81, 122, 117, 97, 97, 0, 121, 96, 121, 0, 121, 88, 128, 97, 82, 121, 0, 121, 130, 119, 0, 81, 121, 2, 253, 0, 159, 99, 97, 81, 121, 81, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 81, 122, 117, 103, 81, 121, 3, 0, 0, 1, 159, 99, 1, 253, 97, 82, 121, 82, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 103, 81, 121, 5, 0, 0, 0, 0, 1, 159, 99, 1, 254, 97, 82, 121, 84, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 103, 81, 121, 9, 0, 0, 0, 0, 0, 0, 0, 0, 1, 159, 99, 1, 255, 97, 82, 121, 88, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 104, 104, 104, 104, 0, 121, 83, 121, 126, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 81, 122, 117, 97, 82, 121, 81, 121, 126, 83, 122, 117, 82, 122, 82, 122, 82, 121, 117, 117, 117, 104, 92, 121, 0, 160, 99, 97, 93, 121, 90, 121, 89, 121, 126, 1, 20, 126, 81, 121, 126, 90, 121, 126, 88, 121, 126, 81, 122, 117, 97, 97, 0, 121, 94, 121, 0, 121, 88, 128, 97, 82, 121, 0, 121, 130, 119, 0, 81, 121, 2, 253, 0, 159, 99, 97, 81, 121, 81, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 81, 122, 117, 103, 81, 121, 3, 0, 0, 1, 159, 99, 1, 253, 97, 82, 121, 82, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 103, 81, 121, 5, 0, 0, 0, 0, 1, 159, 99, 1, 254, 97, 82, 121, 84, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 103, 81, 121, 9, 0, 0, 0, 0, 0, 0, 0, 0, 1, 159, 99, 1, 255, 97, 82, 121, 88, 81, 121, 81, 121, 81, 147, 128, 0, 121, 81, 121, 130, 119, 81, 148, 127, 117, 0, 127, 119, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 104, 104, 104, 104, 0, 121, 83, 121, 126, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 126, 81, 122, 117, 81, 122, 117, 97, 82, 121, 81, 121, 126, 83, 122, 117, 82, 122, 82, 122, 82, 121, 117, 117, 117, 104, 0, 121, 170, 0, 121, 97, 1, 22, 121, 0, 121, 130, 119, 81, 121, 81, 121, 88, 148, 127, 117, 81, 121, 1, 40, 148, 127, 119, 81, 122, 117, 81, 122, 117, 97, 135, 105, 1, 22, 121, 169, 84, 121, 135, 105, 1, 23, 121, 1, 23, 121, 172, 105, 97, 1, 21, 121, 97, 0, 121, 32, 151, 223, 215, 104, 81, 191, 70, 94, 143, 113, 85, 147, 178, 23, 113, 72, 88, 187, 233, 87, 15, 243, 189, 94, 51, 132, 10, 52, 226, 15, 240, 38, 33, 2, 186, 121, 223, 95, 138, 231, 96, 74, 152, 48, 240, 60, 121, 51, 2, 129, 134, 174, 222, 6, 117, 161, 111, 2, 93, 196, 248, 190, 142, 236, 3, 130, 33, 10, 196, 7, 240, 228, 189, 68, 191, 194, 7, 53, 90, 119, 139, 4, 98, 37, 167, 6, 143, 197, 158, 231, 237, 164, 58, 217, 5, 170, 219, 255, 200, 0, 32, 108, 38, 107, 48, 230, 161, 49, 156, 102, 220, 64, 30, 91, 214, 180, 50, 186, 73, 104, 142, 236, 209, 24, 41, 112, 65, 218, 128, 116, 206, 8, 16, 32, 16, 8, 206, 116, 128, 218, 65, 112, 41, 24, 209, 236, 142, 104, 73, 186, 50, 180, 214, 91, 30, 64, 220, 102, 156, 49, 161, 230, 48, 107, 38, 108, 1, 19, 121, 1, 19, 121, 133, 86, 121, 170, 97, 97, 0, 121, 0, 121, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 1, 0, 126, 129, 81, 122, 117, 97, 87, 121, 86, 121, 86, 121, 86, 121, 86, 121, 83, 121, 86, 121, 84, 121, 87, 121, 149, 147, 149, 33, 65, 65, 54, 208, 140, 94, 210, 191, 59, 160, 72, 175, 230, 220, 174, 186, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 97, 81, 121, 81, 121, 81, 121, 81, 121, 151, 82, 122, 117, 81, 122, 81, 121, 0, 159, 99, 81, 121, 81, 121, 147, 82, 122, 117, 81, 122, 104, 81, 121, 81, 122, 117, 81, 122, 117, 97, 82, 122, 117, 81, 122, 81, 121, 81, 121, 82, 150, 160, 99, 0, 121, 82, 121, 148, 82, 122, 117, 81, 122, 104, 83, 121, 130, 119, 82, 121, 130, 119, 84, 82, 121, 147, 81, 121, 147, 1, 48, 81, 121, 126, 82, 126, 83, 121, 126, 87, 121, 126, 82, 126, 82, 121, 126, 85, 121, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 81, 127, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 124, 126, 126, 86, 121, 126, 0, 121, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 0, 121, 87, 121, 172, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 81, 122, 117, 97, 81, 122, 117, 97, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 106}

// These are used for the accepted coin types in regard to wallet functions
const (
	CoinTypeBSV coinType = "BSV" // BitcoinSV
	CoinTypeBTC coinType = "BTC" // BitcoinCore
	CoinTypeETH coinType = "ETH" // Ethereum
)

// coinType is used for determining the coin_type for wallet functions
type coinType string

// String is the string version of coin_type
func (c coinType) String() string {
	return string(c)
}

// StandardResponse is the standard fields returned on all responses from Request()
type StandardResponse struct {
	Body       []byte          `json:"-"` // Body of the response request
	Error      *Error          `json:"-"` // API error response
	StatusCode int             `json:"-"` // Status code returned on the request
	Tracing    resty.TraceInfo `json:"-"` // Trace information if enabled on the request
}

// genericResponse is the generic part of the response body
type genericResponse struct {
	Code    int    `json:"code"` // If there is an error, this will be a value
	Message string `json:"msg"`  // If there is an error, this will be the error message
}

/*
Example
{
    "code": 74012,
    "msg": "Client authentication failed",
    "data": null,
    "req_id": "zacd8b1b05b12a36d45fvfc20a4b97c5"
}
*/

// Error is the universal Error response from the API
//
// For more information: https://developers.dotwallet.com/documents/en/#errors
type Error struct {
	genericResponse
	Data      interface{} `json:"data"`
	Method    string      `json:"method"`
	RequestID string      `json:"req_id"`
	URL       string      `json:"url"`
}

// getAccessTokenRequest is used for the access token request
//
// For more information: https://developers.dotwallet.com/documents/en/#application-authorization
type getAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

// getDotUserTokenRequest is used for the GetDotUserToken request
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
type getDotUserTokenRequest struct {
	getAccessTokenRequest
	Code        string `json:"code"`         // User code given from the oauth2 handshake
	RedirectURI string `json:"redirect_uri"` // The redirect URI set for the application
}

// refreshDotUserTokenRequest is used for the RefreshDotUserToken request
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
type refreshDotUserTokenRequest struct {
	getAccessTokenRequest
	RefreshToken string `json:"refresh_token"` // Refresh token which was given upon first auth_token generation
}

// DotAccessToken is the access token information
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
type DotAccessToken struct {
	AccessToken           string   `json:"access_token"`                       // Access token from the API
	ExpiresAt             int64    `json:"expires_at,omitempty"`               // Friendly unix time from UTC when the access_token expires
	ExpiresIn             int64    `json:"expires_in"`                         // Seconds from now that the token expires
	RefreshToken          string   `json:"refresh_token,omitempty"`            // Refresh token for the user
	RefreshTokenExpiresIn int64    `json:"refresh_token_expires_in,omitempty"` // Seconds from now that the token expires
	RefreshTokenExpiresAt int64    `json:"refresh_token_expires_at,omitempty"` // Friendly unix time from UTC when the refresh_token expires
	Scopes                []string `json:"scopes,omitempty"`                   // Scopes for the user token
	TokenType             string   `json:"token_type"`                         // Token type
}

// accessTokenResponse is the response from creating the new access token
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
type accessTokenResponse struct { //nolint: musttag // This struct was not created properly following Go conventions
	genericResponse
	Data struct {
		AccessToken  string `json:"access_token"`            // Access token from the API
		ExpiresIn    int64  `json:"expires_in"`              // Seconds from now that the token expires
		RefreshToken string `json:"refresh_token,omitempty"` // Refresh token for the user
		Scope        string `json:"scope,omitempty"`         // Scopes for the user token
		TokenType    string `json:"token_type"`              // Token type
	}
}

// userResponse is the response from the user info request
//
// For more information: https://developers.dotwallet.com/documents/en/#user-info
type userResponse struct { //nolint: musttag // This struct was not created properly following Go conventions
	genericResponse
	Data struct {
		User
	}
}

type nftResponse struct { //nolint: musttag // This struct was not created properly following Go conventions
	genericResponse
	Data struct {
		NftData
	}
}

type nftMintResponse struct { //nolint: musttag // This struct was not created properly following Go conventions
	genericResponse
	Data struct {
		NftMintData
	}
}

type transferNftToAddressResponse struct { //nolint: musttag // This struct was not created properly following Go conventions
	genericResponse
	Data struct {
		TransferNftToAddressData
	}
}

// User is the DotWallet user profile information
//
// For more information: https://developers.dotwallet.com/documents/en/#user-info
type User struct {
	Avatar           string            `json:"avatar"`
	ID               string            `json:"id"`
	Nickname         string            `json:"nickname"`
	WebWalletAddress *webWalletAddress `json:"web_wallet_address"`
}

// NftData is the NFT data
type NftData struct {
	CodeHash string `json:"code_hash"`
	Param    string `json:"param"`
}

// NftMintData is the data struct
type NftMintData struct {
	Fee       int64    `json:"fee"`
	FeeStr    string   `json:"fee_str"`
	TxID      string   `json:"txid"`
	BadgePath []string `json:"badge_path"`
}

// TransferNftToAddressData is the transfer to address struct
type TransferNftToAddressData struct {
	Fee    int64  `json:"fee"`
	FeeStr string `json:"fee_str"`
}

// webWalletAddress is the user's wallet addresses
type webWalletAddress struct {
	BSVRegular string `json:"bsv_regular"`
	BTCRegular string `json:"btc_regular"`
	ETHRegular string `json:"eth_regular"`
}

// userReceiveRequest is used for the user receive address request
//
// For more information: https://developers.dotwallet.com/documents/en/#get-user-receive-address
type userReceiveRequest struct {
	UserID   string   `json:"user_id"`
	CoinType coinType `json:"coin_type"`
}

type getNftParam struct {
	TxID string `json:"txid"`
}

type mintNftParam struct {
	CodeHash string `json:"code_hash"`
	Param    string `json:"param"`
}

type transferNftToAddressParam struct {
	TxID    string `json:"txid"`
	Address string `json:"address"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	PicURL  string `json:"pic_url"`
}

// userReceiveAddressResponse is the response from the user receive address request
//
// For more information: https://developers.dotwallet.com/documents/en/#get-user-receive-address
type userReceiveAddressResponse struct { //nolint: musttag // This struct was not created properly following Go conventions
	genericResponse
	Data struct {
		Wallets
	}
}

// Wallets is the user's wallet information
type Wallets struct {
	AutopayWallet *walletInfo `json:"autopay_wallet"`
	PrimaryWallet *walletInfo `json:"primary_wallet"`
}

// walletInfo is the user's wallet information
type walletInfo struct {
	Address     string `json:"address"`
	CoinType    string `json:"coin_type"`
	Paymail     string `json:"paymail,omitempty"`
	UserID      string `json:"user_id"`
	WalletIndex int64  `json:"wallet_index"`
	WalletType  string `json:"wallet_type"`
}

// Target is the target
type Target struct {
	Hash              string  `json:"hash"`
	Confirmations     int     `json:"confirmations"`
	Height            int     `json:"height"`
	Version           int     `json:"version"`
	VersionHex        string  `json:"versionHex"`
	Merkleroot        string  `json:"merkleroot"`
	NumTx             int     `json:"num_tx"`
	Time              int     `json:"time"`
	Mediantime        int     `json:"mediantime"`
	Nonce             int     `json:"nonce"`
	Bits              string  `json:"bits"`
	Difficulty        float64 `json:"difficulty"`
	Chainwork         string  `json:"chainwork"`
	Previousblockhash string  `json:"previousblockhash"`
	Nextblockhash     string  `json:"nextblockhash"`
}

// MerkelProof is the merkle proof
type MerkelProof struct {
	Flags  int      `json:"flags"`
	Index  int      `json:"index"`
	Nodes  []string `json:"nodes"`
	Target *Target  `json:"target"`
	TxOrID string   `json:"txOrId"`
}

// GetMerkleProofResponse is the response
type GetMerkleProofResponse struct { //nolint: musttag // This struct was not created properly following Go conventions
	genericResponse
	Data struct {
		MerkelProof
	}
}

// GetMerkleProofRequest is the request
type GetMerkleProofRequest struct {
	Txid string `json:"txid"`
}

// NftAuthInfo is the auth info
type NftAuthInfo struct {
	Sig []byte `json:"sig"`
	Pub []byte `json:"pub"`
}

// MsgTxInfo is the msg tx info
type MsgTxInfo struct {
	MsgTx     *wire.MsgTx
	Height    int64
	BlockHash string
	Timestamp int64
	RawTx     string
}

// GetRawTransactionResponse is the response
type GetRawTransactionResponse struct { //nolint: musttag // This struct was not created properly following Go conventions
	genericResponse
	Data struct {
		MsgTxInfo
	}
}

// GetRawTransactionRequest is the request
type GetRawTransactionRequest struct {
	TxID string `json:"txid"`
}

// NFT Types
const (
	NftTxTypeCasting    = 1
	NftTxTypeTransfer   = 2
	NftTxTypeDestroy    = 3
	NftTxTypeError      = 4
	NftTxTypeIrrelevant = 5
)

// TxIDIndexPair is the pair
type TxIDIndexPair struct {
	TxID  string
	Index int
}

// AddressBadgeCodePair is the badge code pair
type AddressBadgeCodePair struct {
	Address   string `json:"address"`
	BadgeCode string `json:"badge_code"`
}

// AddressIndexPair is the index pair
type AddressIndexPair struct {
	Address string
	Index   int
}

// NftTxInfo is the tx info
type NftTxInfo struct {
	TxID            string
	NftPreOutPoints []*TxIDIndexPair
	NftOutPoints    []*AddressIndexPair
	Type            int
}
