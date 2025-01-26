package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/client/v3"
	"websocket_proxy/options"
)

type EtcdClient struct {
	*clientv3.Client
}

func NewEtcdClient(opts *options.EtcdOptions) IRegistry {
	config := clientv3.Config{
		Endpoints:   opts.EtcdEndPoints, // etcd 集群地址
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		log.Fatal("Failed to create etcd client: ", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = client.Status(ctx, config.Endpoints[0])
	if err != nil {
		log.Fatal("Failed to connect to etcd: ", err)
		return nil
	}

	_, err = client.Get(ctx, "health_check")
	if err != nil {
		log.Fatal("Failed to connect to etcd: ", err)
		return nil
	}

	return &EtcdClient{Client: client}
}

func (etcd *EtcdClient) GetType() string {
	return options.ETCD
}

func (etcd *EtcdClient) PutServer(prefix string, info ServerInfo, ttl int) error {
	key := fmt.Sprintf("%s%d", prefix, info.ServerID)

	leaseResp, err := etcd.Grant(context.Background(), int64(ttl))
	if err != nil {
		return fmt.Errorf("failed to create lease: %v", err)
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = etcd.Put(context.Background(), key, string(data), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *EtcdClient) GetAllServer(prefix string) map[string]ServerInfo {
	resp, err := etcd.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Printf("failed to get data from etcd with prefix %s: %v", prefix, err)
		return nil
	}

	data := make(map[string]ServerInfo)
	for _, kv := range resp.Kvs {
		info := &ServerInfo{}
		if err = json.Unmarshal(kv.Value, info); err != nil {
			log.Printf("failed to unmarshal server info: %v", err)
			return nil
		}
		data[string(kv.Key)] = *info
	}
	return data
}

func (etcd *EtcdClient) DeleteServer(prefix string, serverID int) error {
	key := fmt.Sprintf("%s%d", prefix, serverID)
	_, err := etcd.Delete(context.Background(), key)
	if err != nil {
		return fmt.Errorf("failed to delete server from etcd: %v", err)
	}
	return nil
}

func (etcd *EtcdClient) Close() {
	etcd.Client.Close()
}
