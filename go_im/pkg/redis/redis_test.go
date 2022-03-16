package redis

import (
	"testing"
)

// config=./config.yaml go test -v

func TestGetInstance(t *testing.T) {
	r := GetInstance()
	t.Log(r.Do(`ping`))
}
