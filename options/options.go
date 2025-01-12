package options

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Options struct {
	Servers        []ServerOptions `yaml:"Servers"`
	Token          *TokenOptions   `yaml:"Token"`
	Secure         *SecureOptions  `yaml:"Secure"`
	Etcd           *EtcdOptions    `yaml:"Etcd"`
	Redis          *RedisOptions   `yaml:"Redis"`
	RegistrySelect string          `yaml:"RegistrySelect"`
}

func Load(filePath string) (*Options, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var opts Options
	err = yaml.Unmarshal(fileContent, &opts)
	if err != nil {
		return nil, err
	}

	return &opts, nil
}

func (opts *Options) GetServerOptions(serverID int) *ServerOptions {
	for _, server := range opts.Servers {
		if server.ServerID == serverID {
			return &server
		}
	}
	log.Fatalf("Server options not found, serverID:%d", serverID)
	return nil
}

func (opts *Options) GetRegistryOptions() IRegistryOptions {
	switch opts.RegistrySelect {
	case "ETCD":
		return opts.Etcd
	case "REDIS":
		return opts.Redis
	default:
		log.Fatalln("No Registry options selected")
		return nil
	}
}

type IRegistryOptions interface {
	GetAddr() string
	GetAddrs() []string
	GetPassword() string
	GetKeyPrefix() string
	GetKeyExpireTime() int
	GetUpdateDuration() int
}

type RegistryOptions struct {
	Key            string `yaml:"Key"`
	Password       string `yaml:"Password"`
	KeyExpireTime  int    `yaml:"KeyExpireTime"`  // 节点过期时间
	UpdateDuration int    `yaml:"UpdateDuration"` // 更新节点间隔
}

func (r RegistryOptions) GetKeyPrefix() string {
	return r.Key
}

func (r RegistryOptions) GetPassword() string {
	return r.Password
}

func (r RegistryOptions) GetKeyExpireTime() int {
	return r.KeyExpireTime
}

func (r RegistryOptions) GetUpdateDuration() int {
	return r.UpdateDuration
}

func (r RegistryOptions) GetAddr() string {
	return ""
}

func (r RegistryOptions) GetAddrs() []string {
	return []string{}
}

func (r RegistryOptions) String() string {
	return fmt.Sprintf("Key:%s,Password:%s,KeyExpireTime:%d,UpdateDuration:%d",
		r.Key, r.Password, r.KeyExpireTime, r.UpdateDuration)
}
