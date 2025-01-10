package options

import "fmt"

type EtcdOptions struct {
	EtcdEndPoints      []string `yaml:"EtcdEndPoints"`
	EtcdKey            string   `yaml:"EtcdKey"`
	EtcdLeaseTime      int      `yaml:"EtcdLeaseTime"`      // 节点过期时间
	UpdateEtcdDuration int      `yaml:"UpdateEtcdDuration"` // 更新节点间隔
}

func NewEtcdOptions() *EtcdOptions {
	return &EtcdOptions{}
}

func (opts *EtcdOptions) String() string {
	return fmt.Sprintf("EtcdEndPoints:%s, EtcdKey:%s, EtcdLeaseTime:%d, UpdateEtcdDuration:%d",
		opts.EtcdEndPoints, opts.EtcdKey, opts.EtcdLeaseTime, opts.UpdateEtcdDuration)
}
