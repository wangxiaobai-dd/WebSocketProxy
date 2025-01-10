package options

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Options struct {
	Servers []ServerOptions `yaml:"Servers"`
	Token   *TokenOptions   `yaml:"Token"`
	SSL     *SSLOptions     `yaml:"SSL"`
	Etcd    *EtcdOptions    `yaml:"Etcd"`
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
	return nil
}
