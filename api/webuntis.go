package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

type webuntis struct {
	Url        string
	Schoolname string
	HttpC      *http.Client
}

func New(url string, schoolname string) webuntis {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	w := webuntis{
		Url:        url,
		Schoolname: schoolname,
		HttpC:      client,
	}
	return w
}
func (w webuntis) Authenticate() (err error) {
	params := map[string]any{}
	return
}

func (w webuntis) Logout() (err error) {

	return
}

func (w webuntis) sendRequest(method string, params map[string]any) (response map[string]any, err error) {
	body := map[string]any{
		"id":      "test",
		"method":  method,
		"jsonrpc": "2.0",
		"params":  params,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, w.Url, bytes.NewReader(bodyJson))
	if err != nil {
		return
	}
	resp, err := w.HttpC.Do(req)

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(responseBody, &response)
	return
}
