package options

import "fmt"

type TokenOptions struct {
	TokenValidTime     int `json:"TokenValidTime" yaml:"TokenValidTime"` // token有效时间
	CheckTokenDuration int `json:"CheckTokenDuration" yaml:"CheckTokenDuration"`
}

func NewTokenOptions() *TokenOptions {
	return &TokenOptions{}
}

func (opts *TokenOptions) String() string {
	return fmt.Sprintf("TokenValidTime:%d,CheckTokenDuration:%d", opts.TokenValidTime, opts.CheckTokenDuration)
}
