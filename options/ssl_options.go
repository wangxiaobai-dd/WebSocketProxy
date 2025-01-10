package options

import "fmt"

type SSLOptions struct {
	Secure   bool   `json:"Secure" yaml:"Secure"`
	CertFile string `json:"CertFile" yaml:"CertFile"`
	KeyFile  string `json:"KeyFile" yaml:"KeyFile"`
}

func NewSSLOptions() *SSLOptions {
	return &SSLOptions{}
}

func (opts *SSLOptions) String() string {
	return fmt.Sprintf("Secure:%v, CertFile:%s, KeyFile:%s", opts.Secure, opts.CertFile, opts.KeyFile)
}
