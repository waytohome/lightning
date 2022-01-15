package promx

import (
	"math/rand"
	"testing"
	"time"

	"github.com/waytohome/lightning/logx"
)

func TestCollectServerRequest(t *testing.T) {
	rand.Seed(time.Now().Unix())
	now := time.Now()
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	CollectServerRequest("unit-test", "GET", "/hello", "200", "unknown", now)

	time.Sleep(20 * time.Second)
}

func TestCollectMySQLRequest(t *testing.T) {
	rand.Seed(time.Now().Unix())
	now := time.Now()
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	CollectMySQLRequest("test", "test", "select", now)

	time.Sleep(20 * time.Second)
}

func init() {
	logx.SetLevel("info")
	StartPusher("192.168.31.13", "9091", "local-unit-test", 15)
}
