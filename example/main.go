package main

import (
	"dotwallet"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var htmlStr string = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>

<body>
    <div id="displayZoom">
        <button type="button" onclick="dotWalletAuth()">dotwallet Login</button>
        <br>
        <button type="button" hidden="hidden" id="GetUserReceiveAddressParaList" onclick="ShowGetUserReceiveAddress()">
            get user receive address
        </button>


        <div hidden="hidden" id="GetUserReceiveAddress">
            <span>coin type: </span>
            <input type="text" value="BSV" id="GetUserReceiveAddress_coin_type">
            <br>
            <button type="button" onclick="GetUserReceiveAddress()">
                submit:
            </button>
            <div>result</div>
            <div style="border:1px solid;height:100px" id="GetUserReceiveAddressResult"></div>
        </div>
        <br>



        <button type="button" hidden="hidden" id="GetAutoPayParamList" onclick="ShowAutoPay()">
            auto pay
        </button>

        <div hidden="hidden" id="AutoPay">
            <span>coin type: </span>
            <input type="text" value="BSV" id="AutoPay_coin_type">
            <br>

            <span>subject: </span>
            <input type="text" value="subject" id="AutoPay_subject">
            <br>

            <span>to: </span>
            <br>

            <span>amount: </span>
            <input type="text" value="1000" id="AutoPay_amount">


            <span>content: </span>
            <input type="text" value="n36EQf9dMkVMhSoMCW1eWHgs6L5xAoJZT9" id="AutoPay_content">

            <span>type: </span>
            <input type="text" value="address" id="AutoPay_type">
            <br>


            <span>product: </span>
            <br>



            <span>id: </span>
            <input type="text" value="id" id="AutoPay_id">


            <span>name: </span>
            <input type="text" value="name" id="AutoPay_name">

            <span>detail: </span>
            <input type="text" value="detail" id="AutoPay_detail">
            <br>


            <br>


            <button type="button" onclick="AutoPay()">
                submit:
            </button>
            <div>result</div>
            <div style="border:1px solid;height:100px" id="AutoPayResult"></div>
        </div>

    </div>
    <script>

        var id = "";

        function DoAjax(method, uri, header, reqdata, callbackWhenReadyState4) {
            request = new XMLHttpRequest();
            request.onreadystatechange = function () {
                if (request.readyState == 4 && request.status == 200) {
                    return callbackWhenReadyState4(request)
                }
            }
            request.open(method, uri, true)
            request.send(JSON.stringify(reqdata))
        }



        function getQuery(key) {
            var query = window.location.search.substring(1);
            var map = query.split("&");
            for (var i = 0; i < map.length; i++) {
                var pair = map[i].split("=");
                if (pair[0] == key) {
                    return pair[1];
                }
            }
        }
        function printAjaxResult(request) {
            if (request.readyState == 4 && request.status == 200) {
                console.log(request.responseText)
            }
        }
        function passwordLogin() {
            let username = document.getElementById("username").value;

            let password = document.getElementById("password").value;


            let UsernamePassword = Object();

            UsernamePassword.username = username
            UsernamePassword.password = password
            DoAjax("POST", "auth", null, UsernamePassword, function (request) {
                console.log(request.responseText)
                let response = JSON.parse(request.responseText)
                if (response.code != 0) {
                    alert(response.msg)
                    return
                }
                id = response.data.id
                console.log("id:", id)
                alert("log in successful")
            })
        }
        function dotWalletAuth() {
            DoAjax("GET", "dot_wallet_auth", null, null, function () {
                window.location.replace(request.responseText)
                return
            })
        }




        function dotWalletLogin() {
            let code = decodeURI(getQuery("code"))
            let state = decodeURI(getQuery("state"))
            if (code != "undefined" && state != "undefined") {
                let CodeState = Object();
                CodeState.code = code
                CodeState.state = state
                DoAjax("POST", "dot_wallet_Login", null, CodeState, function (request) {
                    let response = JSON.parse(request.responseText)
                    if (response.code != 0) {
                        alert(response.msg)
                        return
                    }
                    id = response.data.id
                    let GetUserReceiveAddressParaList = document.getElementById("GetUserReceiveAddressParaList")
                    GetUserReceiveAddressParaList.removeAttribute("hidden")
                    let GetAutoPayParamList = document.getElementById("GetAutoPayParamList")
                    GetAutoPayParamList.removeAttribute("hidden")
                })
            }
        }

        function Hide(elemId) {
            let elem = document.getElementById(elemId)
            elem.removeAttribute("hidden")
        }


        function DisplayOrHide(elemId) {
            let elem = document.getElementById(elemId)
            let hidden = elem.getAttribute("hidden")
            if (hidden == "hidden") {
                elem.removeAttribute("hidden")
                return
            }
            elem.setAttribute("hidden", "hidden")
        }


        function ShowGetUserReceiveAddress() {
            // Hide("AutoPay")
            DisplayOrHide("GetUserReceiveAddress")
        }

        function ShowAutoPay() {
            // Hide("GetUserReceiveAddress")
            DisplayOrHide("AutoPay")
        }


        function GetUserReceiveAddress() {
            let coinType = document.getElementById("GetUserReceiveAddress_coin_type").value;

            let GetUserReceiveAddressRequest = Object();
            GetUserReceiveAddressRequest.coin_type = coinType
            GetUserReceiveAddressRequest.id = id
            DoAjax("POST", "get_user_receive_address", null, GetUserReceiveAddressRequest, function () {
                let GetUserReceiveAddressResult = document.getElementById("GetUserReceiveAddressResult")
                GetUserReceiveAddressResult.innerHTML = request.responseText
            })
        }

        function AutoPay() {
            let amountStr = document.getElementById("AutoPay_amount").value;
            let amount = parseInt(amountStr);
            let content = document.getElementById("AutoPay_content").value;
            let type = document.getElementById("AutoPay_type").value;


            let to = Array();

            let toPoint = Object()

            toPoint.content = content
            toPoint.amount = amount
            toPoint.type = type
            to[0] = toPoint

            let productid = document.getElementById("AutoPay_id").value;
            let name = document.getElementById("AutoPay_name").value;
            let detail = document.getElementById("AutoPay_detail").value;


            let product = Object();
            product.id = productid
            product.name = name
            product.detail = detail


            let coinType = document.getElementById("AutoPay_coin_type").value;

            let subject = document.getElementById("AutoPay_subject").value;


            let AutoPayRequest = Object();
            AutoPayRequest.coin_type = coinType
            AutoPayRequest.user_id = id
            AutoPayRequest.subject = subject
            AutoPayRequest.to = to
            AutoPayRequest.product = product

            DoAjax("POST", "auto_pay", null, AutoPayRequest, function () {
                let AutoPayResult = document.getElementById("AutoPayResult")
                AutoPayResult.innerHTML = request.responseText
            })
        }



        dotWalletLogin()




    </script>
</body>

</html>`

var gClient *dotwallet.Client

var gStates = make(map[string]bool)

var gConfig = &Config{}

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
	state := uuid.NewV4().String()
	gStates[state] = true
	rsp.Write(
		[]byte(
			gClient.GetAuthorizeUrl(
				state,
				[]string{
					dotwallet.SCOPE_USER_INFO,
					dotwallet.SCOPE_AUTOPAY_BSV,
					dotwallet.SCOPE_AUTOPAY_BTC,
					dotwallet.SCOPE_AUTOPAY_ETH,
				},
			),
		),
	)
}

type DotWalletLoginRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type LoginResponse struct {
	Id string `json:"id"`
}

func DotWalletLogin(rsp http.ResponseWriter, req *http.Request) {
	dotWalletLoginRequest := &DotWalletLoginRequest{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	err = json.Unmarshal(body, dotWalletLoginRequest)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	_, ok := gStates[dotWalletLoginRequest.State]
	if !ok {
		rsp.Write(MakeErrHttpJsonResponse(-1, "state not found"))
		return
	}
	dotUser, err := gClient.GetDotUser(dotWalletLoginRequest.Code, dotWalletLoginRequest.State)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	rsp.Write(MakeOKHttpJsonResponseByInterface(
		&LoginResponse{
			DotUserId2UserId(dotUser.Id),
		},
	))
}

type GetUserReceiveAddressRequest struct {
	Id       string `json:"id"`
	CoinType string `json:"coin_type"`
}

func GetUserReceiveAddress(rsp http.ResponseWriter, req *http.Request) {
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
		NotifyUrl:  gConfig.NotifyUrl,
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
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rsp.Write(MakeErrHttpJsonResponse(-1, err.Error()))
		return
	}
	fmt.Println(string(body))
}

func LoginPage(rsp http.ResponseWriter, req *http.Request) {
	rsp.Header().Set("Content-Type", "text/html")
	rsp.Write([]byte(htmlStr))
}

func StartHttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/dot_wallet_auth", DotWalletAuth)
	r.HandleFunc("/dot_wallet_Login", DotWalletLogin)
	r.HandleFunc("/get_user_receive_address", GetUserReceiveAddress)
	r.HandleFunc("/auto_pay", AutoPay)
	r.HandleFunc("/auto_pay_notify", AutoPayNotify)
	r.HandleFunc("/Login", LoginPage)
	r.HandleFunc("/", LoginPage)
	err := http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Host         string
	ClientId     string
	ClientSecret string
	RedirectUri  string
	NotifyUrl    string
}

func main() {
	configFilePath := flag.String("config", "./config.json", "Path of config file")
	flag.Parse()
	configJSON, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configJSON, gConfig)
	if err != nil {
		panic(err)
	}

	gClient, err = dotwallet.NewClient(
		gConfig.Host,
		gConfig.ClientId,
		gConfig.ClientSecret,
		gConfig.RedirectUri,
	)
	if err != nil {
		panic(err)
	}

	StartHttpServer()
}
