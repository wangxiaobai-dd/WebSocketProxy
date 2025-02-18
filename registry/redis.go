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
	stopChan chan struct{}
}

func NewRedisClient(opts *options.RedisOptions) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("Redis connect success")
			return nil
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Fatal("Failed to connect Redis client: ", err)
		return nil
	}

	redisClient := &RedisClient{
		Client:   client,
		stopChan: make(chan struct{}),
	}

	if opts.KeepAlive > 0 {
		go func(stopChan chan struct{}) {
			ticker := time.NewTicker(time.Duration(opts.KeepAlive) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err := client.Ping(context.Background()).Err(); err != nil {
						log.Printf("Redis keep-alive failed, err:%v", err)
					}
				case <-stopChan:
					log.Printf("Redis keep-alive stopped")
					return
				}
			}
		}(redisClient.stopChan)
	}
	return redisClient
}

func (r *RedisClient) GetType() string {
	return options.REDIS
}

func (r *RedisClient) PutServer(prefix string, info ServerInfo, ttl int) error {
	key := getServerKey(prefix, info.ServerID)
	err := r.ZAdd(context.Background(), getServerConnSetKey(prefix), redis.Z{Score: float64(info.ConnNum), Member: key}).Err()
	if err != nil {
		return err
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
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
	err := r.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("delete server, failed to ping, serverID:%d, err:%v", serverID, err)
	}
	key := getServerKey(prefix, serverID)

	err = r.ZRem(context.Background(), getServerConnSetKey(prefix), serverID).Err()
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
	close(r.stopChan)
	r.Client.Close()
}
