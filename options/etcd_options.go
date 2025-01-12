package options

import "fmt"

type EtcdOptions struct {
	RegistryOptions `yaml:"Registry"`
	EtcdEndPoints   []string `yaml:"EtcdEndPoints"`
}

func (opts EtcdOptions) GetAddrs() []string {
	return opts.EtcdEndPoints
}

func (opts EtcdOptions) String() string {
	return fmt.Sprintf("EtcdEndPoints:%s,%s", opts.EtcdEndPoints, opts.RegistryOptions)
}
