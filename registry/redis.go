package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"websocket_proxy/options"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(opts *options.RedisOptions) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect Redis client: ", err)
		return nil
	}
	return &RedisClient{client}
}

func (r *RedisClient) GetType() string {
	return options.REDIS
}

func (r *RedisClient) PutData(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) PutDataWithTTL(key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.SetEx(context.Background(), key, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) GetData(key string, result interface{}) error {
	data, err := r.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		} else {
			return err
		}
	}

	err = json.Unmarshal([]byte(data), result)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetDataWithPrefix(prefix string) (map[string]string, error) {
	ctx := context.Background()
	iter := r.Scan(ctx, 0, prefix+"*", 0).Iterator()

	data := make(map[string]string)
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := r.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get data for key %s, %v", key, err)
		}
		data[key] = val
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan keys with prefix %s, %v", prefix, err)
	}
	return data, nil
}

func (r *RedisClient) DeleteData(key string) error {
	err := r.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) ZAddData(key string, score float64, value interface{}) error {
	err := r.ZAdd(context.Background(), key, redis.Z{Score: score, Member: value}).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) ZRemData(key string, value interface{}) error {
	err := r.ZRem(context.Background(), key, value).Err()
	return err
}

func (r *RedisClient) Close() {
	r.Client.Close()
}
