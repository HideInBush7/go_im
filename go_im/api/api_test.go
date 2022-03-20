package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	resp, err := http.Post(`http://127.0.0.1:8000/user/logout`, `application/json`, strings.NewReader(`{"uid": "3","token":"POl/ly5YTbt0xy8mOM2AfA"}`))
	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("%s\n", body)
}
