package registry

import (
	"log"

	"websocket_proxy/options"
)

type IRegistry interface {
	GetType() string
	PutData(key string, value interface{}) error
	PutDataWithTTL(key string, value interface{}, ttl int) error
	GetData(key string, result interface{}) error
	GetDataWithPrefix(prefix string) (map[string]string, error)
	DeleteData(key string) error
	Close()
}

func NewRegistry(opts options.IRegistryOptions) IRegistry {
	switch v := opts.(type) {
	case *options.EtcdOptions:
		return NewEtcdClient(v)
	case *options.RedisOptions:
		return NewRedisClient(v)
	default:
		log.Fatal("Invalid options: ", v)
		return nil
	}
}
