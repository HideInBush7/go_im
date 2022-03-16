package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	resp, err := http.Post(`http://127.0.0.1:8000/user/register`, `application/json`, strings.NewReader(`{"username": "hello","password":"world"}`))
	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("%s\n", body)
}
