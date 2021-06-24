package dotwallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	HTTP_POST             = "POST"
	HTTP_GET              = "GET"
	HTTP_CONTENT_TYPE     = "Content-Type"
	HTTP_APPLICATION_JSON = "application/json"
)

func ToJson(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func ToCurlStr(method string, header map[string]string, body []byte, url string) {
	var b strings.Builder
	b.WriteString("curl ")
	for key, value := range header {
		b.WriteString("-H ")
		b.WriteString("\"")
		b.WriteString(key)
		b.WriteString(":")
		b.WriteString(value)
		b.WriteString("\" ")
	}
	b.WriteString("-X ")
	b.WriteString(method)

	b.WriteString(" --data '")
	b.Write(body)
	b.WriteString("' ")
	b.WriteString(url)
	fmt.Println(b.String())
}

func DoHttpRequest(method string, url string, urlValues *url.Values, headers map[string]string, reqBody interface{}) ([]byte, error) {
	httpClient := &http.Client{}
	contentByte := make([]byte, 0, 8)
	if reqBody != nil {
		contentByteTmp, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		contentByte = contentByteTmp
	}
	if urlValues != nil {
		url = fmt.Sprintf("%s?%s", url, urlValues.Encode())
	}
	// fmt.Println(url)
	request, err := http.NewRequest(method, url, bytes.NewReader(contentByte))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	ToCurlStr(method, headers, contentByte, url)
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
