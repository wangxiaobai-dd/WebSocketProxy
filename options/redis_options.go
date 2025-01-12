package options

import "fmt"

type RedisOptions struct {
	RegistryOptions `yaml:"Registry"`
	Addr            string `yaml:"Addr"`
}

func (opts RedisOptions) String() string {
	return fmt.Sprintf("Addr:%s,%s", opts.Addr, opts.RegistryOptions)
}
