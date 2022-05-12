package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type webuntis struct {
	Url   string
	HttpC *http.Client
	Id    float64
}

func NewWebUntis(url string, schoolname string) webuntis {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	if url[len(url)-1] != '/' {
		url = url + "/"
	}
	w := webuntis{
		Url:   fmt.Sprintf("%sWebUntis/jsonrpc.do?school=%s", url, strings.Replace(schoolname, " ", "%20", -1)),
		HttpC: client,
	}
	return w
}
func (w *webuntis) Authenticate(username string, password string) (response map[string]any, err error) {
	params := map[string]any{
		"user":     username,
		"password": password,
	}
	response, err = w.sendRequest("authenticate", params)
	if err != nil {
		return
	}
	w.Id = response["result"].(map[string]any)["personId"].(float64)
	return
}

func (w webuntis) Logout() (response map[string]any, err error) {
	response, err = w.sendRequest("logout", map[string]any{})
	return
}

func (w webuntis) GetSubjects() (response map[string]any, err error) {
	response, err = w.sendRequest("getSubjects", map[string]any{})
	return
}

func (w webuntis) GetRooms() (response map[string]any, err error) {
	response, err = w.sendRequest("getRooms", map[string]any{})
	return
}

func (w webuntis) GetTimetableForStudent(daysIntoFuture int) (response map[string]any, err error) {
	today := time.Now().AddDate(0, 0, daysIntoFuture)
	response, err = w.sendRequest("getTimetable", map[string]any{
		"id":        w.Id,
		"type":      5,
		"startDate": today.Format("20060102"),
		"endDate":   today.Format("20060102"),
	})
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
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(responseBody, &response)
	return
}
