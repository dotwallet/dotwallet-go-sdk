package main

import (
	"dotwallet"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	TEST_DOTWALLET_HOST = "http://192.168.1.13:6001"
	TEST_CLIENT_ID      = "06cfc857cfea6028002087e938541d63"
	TEST_CLIENT_SECRET  = "d9fb1a7455bf025462eb75102bcf80ec"
	TEST_REDIRECT_URI   = "http://192.168.1.141:10086/login"
)

var gClient *dotwallet.Client

var gStates map[string]bool = make(map[string]bool)

type HttpJsonResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func MakeHttpJsonResponse(code int, msg string, data json.RawMessage) []byte {
	jsonResponse := &HttpJsonResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	responseByte, err := json.Marshal(jsonResponse)
	if err != nil {
		return []byte(err.Error())
	}
	return responseByte
}

func MakeOKHttpJsonResponse(data json.RawMessage) []byte {
	return MakeHttpJsonResponse(0, "", data)
}

func MakeOKHttpJsonResponseByInterface(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		return MakeErrHttpJsonResponse(-1, err.Error())
	}
	return MakeOKHttpJsonResponse(b)
}

func MakeErrHttpJsonResponse(code int, msg string) []byte {
	return MakeHttpJsonResponse(code, msg, nil)
}

func DotUserId2UserId(dotUserId string) string {
	return dotUserId
}

func UserId2DotUserId(id string) string {
	return id
}

func DotWalletAuth(rsp http.ResponseWriter, req *http.Request) {
	fmt.Println("DotWalletAuth--------")
	state := uuid.NewV4().String()
	fmt.Println("state :", state)
	gStates[state] = true
	rsp.Write(
		[]byte(
			gClient.GetAuthorizeUrl(
				state,
				[]string{
					dotwallet.SCOPT_USER_INFO,
					dotwallet.SCOPT_AUTOPAY_BSV,
					dotwallet.SCOPT_AUTOPAY_BTC,
					dotwallet.SCOPT_AUTOPAY_ETH,
				},
			),
		),
	)
}

type DotWalletLogInRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type LogInResponse struct {
	Id string `json:"id"`
}

func createNewUserReturnId() (string, error) {
	return uuid.NewV4().String(), nil
}

func DotWalletLogIn(rsp http.ResponseWriter, req *http.Request) {
	fmt.Println("DotWalletLogIn----------")
	dotWalletLogInRequest := &DotWalletLogInRequest{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	err = json.Unmarshal(body, dotWalletLogInRequest)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	_, ok := gStates[dotWalletLogInRequest.State]
	if !ok {
		rsp.Write(MakeErrHttpJsonResponse(-1, "state not found"))
		return
	}
	dotUser, err := gClient.GetDotUser(dotWalletLogInRequest.Code, dotWalletLogInRequest.State)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	rsp.Write(MakeOKHttpJsonResponseByInterface(
		&LogInResponse{
			DotUserId2UserId(dotUser.Id),
		},
	))
}

type GetUserReceiveAddressRequest struct {
	Id       string `json:"id"`
	CoinType string `json:"coin_type"`
}

func GetUserReceiveAddress(rsp http.ResponseWriter, req *http.Request) {
	fmt.Println("GetUserReceiveAddress--------")
	getUserReceiveAddressRequest := &GetUserReceiveAddressRequest{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	err = json.Unmarshal(body, getUserReceiveAddressRequest)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	dotId := UserId2DotUserId(getUserReceiveAddressRequest.Id)
	result, err := gClient.GetUserReceiveAddress(dotId, getUserReceiveAddressRequest.CoinType)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	rsp.Write(MakeOKHttpJsonResponseByInterface(
		result,
	))
}

type AutoPayRequest struct {
	CoinType string               `json:"coin_type"`
	UserId   string               `json:"user_id"`
	Subject  string               `json:"subject"`
	Product  *dotwallet.Product   `json:"product"`
	To       []*dotwallet.ToPoint `json:"to"`
}

func AutoPay(rsp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	autoPayRequest := &AutoPayRequest{}
	err = json.Unmarshal(body, autoPayRequest)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	dotAutoPayRequest := &dotwallet.AutoPayRequest{
		OutOrderId: uuid.NewV4().String(),
		CoinType:   autoPayRequest.CoinType,
		UserId:     UserId2DotUserId(autoPayRequest.UserId),
		Subject:    autoPayRequest.Subject,
		NotifyUrl:  "http://192.168.1.141:10086/auto_pay_notify",
		Product:    autoPayRequest.Product,
		To:         autoPayRequest.To,
	}
	result, err := gClient.AutoPay(dotAutoPayRequest)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	rsp.Write(MakeOKHttpJsonResponseByInterface(
		result,
	))
}

func AutoPayNotify(rsp http.ResponseWriter, req *http.Request) {
	fmt.Println("AutoPayNotify-----------")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	fmt.Println(string(body))
}

func LogInPadge(rsp http.ResponseWriter, req *http.Request) {
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rsp.Write([]byte(err.Error()))
		return
	}
	rsp.Header().Set("Content-Type", "text/html")
	html, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	rsp.Write(html)
}

func StartHttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/dot_wallet_auth", DotWalletAuth)
	r.HandleFunc("/dot_wallet_login", DotWalletLogIn)
	r.HandleFunc("/get_user_receive_address", GetUserReceiveAddress)
	r.HandleFunc("/auto_pay", AutoPay)
	r.HandleFunc("/auto_pay_notify", AutoPayNotify)
	r.HandleFunc("/login", LogInPadge)
	r.HandleFunc("/", LogInPadge)
	err := http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		panic(err)
	}
}

func main() {
	client, err := dotwallet.NewClient(
		TEST_DOTWALLET_HOST,
		TEST_CLIENT_ID,
		TEST_CLIENT_SECRET,
		TEST_REDIRECT_URI,
	)
	if err != nil {
		panic(err)
	}
	gClient = client
	StartHttpServer()
}
