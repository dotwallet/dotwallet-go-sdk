package dotwallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	SCOPE_USER_INFO   = "user.info"
	SCOPE_AUTOPAY_BSV = "autopay.bsv"
	SCOPE_AUTOPAY_BTC = "autopay.btc"
	SCOPE_AUTOPAY_ETH = "autopay.eth"

	COIN_BSV = "BSV"
	COIN_BTC = "BTC"
	COIN_ETH = "ETH"

	GET_ACCESS_TOKEN_URI         = "/v1/oauth2/get_access_token"
	GET_USER_INFO_URI            = "/v1/user/get_user_info"
	AUTHORIZE_URI                = "/v1/oauth2/authorize"
	GET_USER_RECEIVE_ADDRESS_URI = "/v1/user/get_user_receive_address"
	AUTOPAY_URI                  = "/v1/transact/order/autopay"

	TO_TYPE_ADDRESS          = "address"
	TO_TYPE_PAYMAIL          = "paymail"
	TO_TYPE_SCRIPT           = "script"
	TO_TYPE_USER_PRIMARY_WEB = "user_primary_web"

	HEADER_AUTHORIZATION = "Authorization"

	GRANT_TYPE_CLIENT_CREDENTIALS = "client_credentials"
	GRANT_TYPE_AUTHORIZATION_CODE = "authorization_code"
	GRANT_TYPE_REFRESH_TOKEN      = "refresh_token"
)

type DotUser struct {
	Id           string
	RefreshToken string
	AccessToken  string
	ExpiredAt    int64
	TokenType    string
	Nickname     string
	Avatar       string
	Scopes       []string
}

type CodeMsgData struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type Client struct {
	token                    string
	tokenType                string
	host                     string
	clientId                 string
	clientSecret             string
	getAccessTokenUrl        string
	authorizeUrl             string
	redirectUri              string
	getUserInfoUrl           string
	getUserReceiveAddressUrl string
	autoPayUrl               string
}

func (this *Client) GetAuthorizeUrl(state string, scopes []string) string {
	urlValues := &url.Values{}
	urlValues.Add("client_id", this.clientId)
	urlValues.Add("redirect_uri", this.redirectUri)
	urlValues.Add("state", state)
	strings.Join(scopes, " ")
	urlValues.Add("scope", strings.Join(scopes, " "))
	urlValues.Add("response_type", "code")
	return fmt.Sprintf("%s%s?%s", this.host, AUTHORIZE_URI, urlValues.Encode())
}

type TokenTokenType struct {
	Token     string
	TokenType string
}

func (this *Client) DoHttpRequestWithToken(
	method string,
	url string,
	urlValues *url.Values,
	headers http.Header,
	reqBody interface{},
	rspData interface{},
) error {
	if headers == nil {
		headers = make(http.Header)
	}
	headers.Set(HEADER_AUTHORIZATION, fmt.Sprintf("%s %s", this.tokenType, this.token))
	err := this.DoHttpRequest(
		method,
		url,
		urlValues,
		headers,
		reqBody,
		rspData,
	)
	if err == nil {
		return nil
	}
	if !strings.Contains(err.Error(), "Your login status has expired") {
		return err
	}
	err = this.UpdateApplicationAccessToken()
	if err != nil {
		return err
	}
	headers.Set(HEADER_AUTHORIZATION, fmt.Sprintf("%s %s", this.tokenType, this.token))
	return this.DoHttpRequest(
		method,
		url,
		urlValues,
		headers,
		reqBody,
		rspData,
	)
}

func (this *Client) DoHttpRequest(
	method string,
	url string,
	urlValues *url.Values,
	headers http.Header,
	reqBody interface{},
	rspData interface{},
) error {
	if headers == nil {
		headers = make(http.Header)
	}
	headers.Set(HTTP_CONTENT_TYPE, HTTP_APPLICATION_JSON)
	body, err := DoHttpRequest(method, url, urlValues, headers, reqBody)
	if err != nil {
		return err
	}
	codeMsgData := &CodeMsgData{}
	err = json.Unmarshal(body, codeMsgData)
	if err != nil {
		return err
	}
	if codeMsgData.Code != 0 {
		return errors.New(codeMsgData.Msg)
	}
	return json.Unmarshal(codeMsgData.Data, rspData)
}

type DotUserTokenInfo struct {
	AccessToken  string   `json:"access_token"`
	ExpiredAt    int64    `json:"expired_at"`
	RefreshToken string   `json:"refresh_token"`
	Scopes       []string `json:"scopes"`
	TokenType    string   `json:"token_type"`
}

type GetDotUserTokenInfoRequest struct {
	ClientId     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
}

type GetAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func NewDotUserTokenInfo(getAccessTokenResponse *GetAccessTokenResponse) *DotUserTokenInfo {
	scopes := strings.Split(getAccessTokenResponse.Scope, " ")
	return &DotUserTokenInfo{
		AccessToken:  getAccessTokenResponse.AccessToken,
		ExpiredAt:    time.Now().Unix() + int64(getAccessTokenResponse.ExpiresIn),
		RefreshToken: getAccessTokenResponse.RefreshToken,
		Scopes:       scopes,
		TokenType:    getAccessTokenResponse.TokenType,
	}
}

func (this *Client) GetUserTokenInfo(code string) (*DotUserTokenInfo, error) {
	getUserTokenInfoRequest := &GetDotUserTokenInfoRequest{
		ClientId:     this.clientId,
		GrantType:    GRANT_TYPE_AUTHORIZATION_CODE,
		ClientSecret: this.clientSecret,
		Code:         code,
		RedirectUri:  this.redirectUri,
	}
	getAccessTokenResponse := &GetAccessTokenResponse{}
	err := this.DoHttpRequest(
		HTTP_POST,
		this.getAccessTokenUrl,
		nil,
		nil,
		getUserTokenInfoRequest,
		getAccessTokenResponse,
	)
	if err != nil {
		return nil, err
	}
	return NewDotUserTokenInfo(getAccessTokenResponse), nil
}

type RefreshTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

func (this *Client) RefreshToken(RefreshToken string) (*DotUserTokenInfo, error) {
	refreshTokenRequest := &RefreshTokenRequest{
		ClientId:     this.clientId,
		ClientSecret: this.clientSecret,
		GrantType:    GRANT_TYPE_REFRESH_TOKEN,
		RefreshToken: RefreshToken,
	}
	getAccessTokenResponse := &GetAccessTokenResponse{}
	err := this.DoHttpRequest(HTTP_POST, this.getAccessTokenUrl, nil, nil, refreshTokenRequest, getAccessTokenResponse)
	if err != nil {
		return nil, err
	}
	return NewDotUserTokenInfo(getAccessTokenResponse), nil
}

type DotUserInfo struct {
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func (this *Client) GetDotUserInfoByUserToken(
	AccessToken string,
	TokenType string,
) (*DotUserInfo, error) {
	header := make(http.Header)
	header.Set(HEADER_AUTHORIZATION, fmt.Sprintf("%s %s", TokenType, AccessToken))
	getUserInfoResponse := &DotUserInfo{}
	err := this.DoHttpRequest(HTTP_POST, this.getUserInfoUrl, nil, header, nil, getUserInfoResponse)
	if err != nil {
		return nil, err
	}
	return getUserInfoResponse, nil
}

func (this *Client) GetDotUser(code string, state string) (*DotUser, error) {
	userTokenInfo, err := this.GetUserTokenInfo(code)
	if err != nil {
		return nil, err
	}
	dotUserInfo, err := this.GetDotUserInfoByUserToken(userTokenInfo.AccessToken, userTokenInfo.TokenType)
	if err != nil {
		return nil, err
	}
	return &DotUser{
		Id:           dotUserInfo.Id,
		AccessToken:  userTokenInfo.AccessToken,
		RefreshToken: userTokenInfo.RefreshToken,
		ExpiredAt:    userTokenInfo.ExpiredAt,
		Scopes:       userTokenInfo.Scopes,
		TokenType:    userTokenInfo.TokenType,
		Avatar:       dotUserInfo.Avatar,
		Nickname:     dotUserInfo.Nickname,
	}, nil
}

type GetApplicationAccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	ClientSecret string `json:"client_secret"`
}

func (this *Client) UpdateApplicationAccessToken() error {
	getAccessTokenRequest := &GetApplicationAccessTokenRequest{
		ClientId:     this.clientId,
		GrantType:    GRANT_TYPE_CLIENT_CREDENTIALS,
		ClientSecret: this.clientSecret,
	}
	getAccessTokenResponse := &GetAccessTokenResponse{}
	err := this.DoHttpRequest(HTTP_POST, this.getAccessTokenUrl, nil, nil, getAccessTokenRequest, getAccessTokenResponse)
	if err != nil {
		return err
	}
	this.token = getAccessTokenResponse.AccessToken
	this.tokenType = getAccessTokenResponse.TokenType
	return nil
}

type GetUserReceiveAddressRequest struct {
	UserId   string `json:"user_id"`
	CoinType string `json:"coin_type"`
}

type UserReceiveAddress struct {
	Address     string `json:"address"`
	Paymail     string `json:"paymail"`
	CoinType    string `json:"coin_type"`
	WalletIndex int64  `json:"wallet_index"`
}

type GetUserReceiveAddressResponse struct {
	PrimaryWallet *UserReceiveAddress `json:"primary_wallet"`
	AutopayWallet *UserReceiveAddress `json:"autopay_wallet"`
}

func (this *Client) GetUserReceiveAddress(id string, coinType string) (*GetUserReceiveAddressResponse, error) {
	getUserReceiveAddressRequest := &GetUserReceiveAddressRequest{
		UserId:   id,
		CoinType: coinType,
	}
	getUserReceiveAddressResponse := &GetUserReceiveAddressResponse{}
	err := this.DoHttpRequestWithToken(
		HTTP_POST,
		this.getUserReceiveAddressUrl,
		nil,
		nil,
		getUserReceiveAddressRequest,
		getUserReceiveAddressResponse,
	)
	if err != nil {
		return nil, err
	}
	return getUserReceiveAddressResponse, nil
}

type ToPoint struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Amount  int64  `json:"amount"`
}

type Product struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

type AutoPayRequest struct {
	OutOrderId string     `json:"out_order_id"`
	CoinType   string     `json:"coin_type"`
	UserId     string     `json:"user_id"`
	Subject    string     `json:"subject"`
	NotifyUrl  string     `json:"notify_url"`
	Product    *Product   `json:"product"`
	To         []*ToPoint `json:"to"`
}

type AutoPayResponse struct {
	OrderId    string `json:"order_id"`
	OutOrderId string `json:"out_order_id"`
	UserId     string `json:"user_id"`
	Amount     int64  `json:"amount"`
	Fee        int64  `json:"fee"`
	Txid       string `json:"txid"`
}

func (this *Client) AutoPay(autoPayRequest *AutoPayRequest) (*AutoPayResponse, error) {
	autoPayResponse := &AutoPayResponse{}
	err := this.DoHttpRequestWithToken(HTTP_POST, this.autoPayUrl, nil, nil, autoPayRequest, autoPayResponse)
	if err != nil {
		return nil, err
	}
	return autoPayResponse, nil
}

func NewClient(
	host string,
	clientId string,
	clientSecret string,
	redirectUri string,
) (*Client, error) {
	client := &Client{
		host:                     host,
		clientId:                 clientId,
		clientSecret:             clientSecret,
		redirectUri:              redirectUri,
		getAccessTokenUrl:        fmt.Sprintf("%s%s", host, GET_ACCESS_TOKEN_URI),
		getUserInfoUrl:           fmt.Sprintf("%s%s", host, GET_USER_INFO_URI),
		getUserReceiveAddressUrl: fmt.Sprintf("%s%s", host, GET_USER_RECEIVE_ADDRESS_URI),
		autoPayUrl:               fmt.Sprintf("%s%s", host, AUTOPAY_URI),
	}
	err := client.UpdateApplicationAccessToken()
	if err != nil {
		return nil, err
	}
	return client, nil
}
