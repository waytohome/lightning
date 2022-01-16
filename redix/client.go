package redix

import "github.com/go-redis/redis/v8"

func NewClient(addr, pwd string) redis.Cmdable {
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
