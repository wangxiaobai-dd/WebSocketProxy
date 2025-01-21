package options

import "fmt"

type SecureOptions struct {
	CertFile string `yaml:"CertFile"`
	KeyFile  string `yaml:"KeyFile"`
}

func (opts SecureOptions) String() string {
	return fmt.Sprintf("CertFile:%s,KeyFile:%s", opts.CertFile, opts.KeyFile)
}
