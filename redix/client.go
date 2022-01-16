package redix

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/waytohome/lightning/confx"
)

func NewClientWithConfigure(c confx.Configure) redis.Cmdable {
	mode, _ := c.GetString("redis.mode", "single")
	addr, _ := c.GetString("redis.addr", "")
	pwd, _ := c.GetString("redis.password", "")
	if mode == "single" {
		return NewSingleClient(addr, pwd)
	} else if mode == "cluster" {
		return NewClusterClient(strings.Split(addr, ","), pwd)
	}
	panic(fmt.Sprintf("config redis.mode %s is unknown", mode))
}

func NewSingleClient(addr, pwd string) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
	})
}

func NewClusterClient(addrs []string, pwd string) redis.Cmdable {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: pwd,
	})
}
