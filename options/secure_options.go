package options

import "fmt"

type SecureOptions struct {
	Enabled  bool   `yaml:"Enabled"`
	CertFile string `yaml:"CertFile"`
	KeyFile  string `yaml:"KeyFile"`
}

func (opts SecureOptions) String() string {
	return fmt.Sprintf("Enabled:%v,CertFile:%s,KeyFile:%s", opts.Enabled, opts.CertFile, opts.KeyFile)
}
