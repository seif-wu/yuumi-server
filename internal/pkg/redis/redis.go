package redis

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	RDB  *redis.Client
	once sync.Once
)

type Option struct {
	Addr     string
	Password string
	DB       int
}

func Init(opts *Option) *redis.Client {
	once.Do(func() {
		RDB = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})

	return RDB
}
