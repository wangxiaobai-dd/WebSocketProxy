package options

import "fmt"

type TokenOptions struct {
	TokenValidTime     int `yaml:"TokenValidTime"` // token有效时间
	CheckTokenDuration int `yaml:"CheckTokenDuration"`
}

func (opts TokenOptions) String() string {
	return fmt.Sprintf("TokenValidTime:%d,CheckTokenDuration:%d", opts.TokenValidTime, opts.CheckTokenDuration)
}
