package registry

import (
	"context"
	"encoding/json"
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

func (r *RedisClient) PutServer(prefix string, info ServerInfo, ttl int) error {
	err := r.ZAdd(context.Background(), getServerConnSetKey(prefix), redis.Z{Score: float64(info.ConnNum), Member: info.ServerID}).Err()
	if err != nil {
		return err
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	key := getServerKey(prefix, info.ServerID)
	err = r.SetEx(context.Background(), key, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetAllServer(prefix string) map[string]ServerInfo {
	ctx := context.Background()
	iter := r.Scan(ctx, 0, prefix+"*", 0).Iterator()

	data := make(map[string]ServerInfo)
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := r.Get(ctx, key).Result()
		if err != nil {
			log.Printf("failed to get data, key:%s, err:%v", key, err)
			continue
		}
		info := &ServerInfo{}
		err = json.Unmarshal([]byte(val), info)
		if err != nil {
			log.Printf("failed to unmarshal data for key %s, %v", key, err)
			continue
		}
		data[key] = *info
	}
	if err := iter.Err(); err != nil {
		log.Printf("failed to scan keys with prefix %s, %v", prefix, err)
	}
	return data
}

func (r *RedisClient) DeleteServer(prefix string, serverID int) error {
	key := getServerKey(prefix, serverID)

	err := r.ZRem(context.Background(), getServerConnSetKey(prefix), serverID).Err()
	if err != nil {
		log.Printf("failed to zrem server prefix:%s, err:%v", prefix, err)
	}
	err = r.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete server prefix:%s, err:%v", prefix, err)
	}
	return nil
}

func (r *RedisClient) Close() {
	r.Client.Close()
}
