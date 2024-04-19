package redis_client

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/koliader/posts-post.git/internal/util"
)

type Client struct {
	client *redis.Client
}

func NewRedis(config util.Config) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: "",
		DB:       config.RedisDBNumber,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	client := Client{
		client: redisClient,
	}
	return &client, nil
}

func (c *Client) Set(key string, value []byte) error {
	err := c.client.Set(context.Background(), key, value, 0).Err()
	return err

}
func (c *Client) Get(key string) (value *string, err error) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return &val, nil
}

func (c *Client) DeleteKey(key string) error {
	err := c.client.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}
