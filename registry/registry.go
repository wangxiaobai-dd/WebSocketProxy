package registry

import (
	"log"

	"websocket_proxy/options"
)

type IRegistry interface {
	GetType() string
	PutServer(prefix string, info ServerInfo, ttl int) error
	GetAllServer(prefix string) map[string]ServerInfo
	DeleteServer(prefix string, serverID int) error
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
