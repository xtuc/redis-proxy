package redis

import (
	"log"

	goredis "github.com/go-redis/redis"

	"github.com/xtuc/redis-proxy/types"
)

type RedisClient struct {
	conn *goredis.Client
}

func NewRedisClient(host string) (types.BackendInterface, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     host + ":6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return &RedisClient{client}, nil
}

func (client RedisClient) GetValue(key string) (string, error) {
	res := client.conn.Get(key)

	if res == nil {
		return "", nil
	}

	return res.Result()
}

func (client RedisClient) PutValue(key, value string) error {
	log.Println(key)
	log.Println(value)

	return nil
}
