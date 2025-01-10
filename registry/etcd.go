package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type EtcdClient struct {
	*clientv3.Client
}

func NewEtcdClient(endPoints []string) *EtcdClient {
	config := clientv3.Config{
		Endpoints:   endPoints, // etcd 集群地址
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		log.Fatal("Failed to create etcd client: ", err)
		return nil
	}
	return &EtcdClient{Client: client}
}

func (etcd *EtcdClient) PutData(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = etcd.Put(context.Background(), key, string(data))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *EtcdClient) GetData(key string, result interface{}) error {
	resp, err := etcd.Get(context.Background(), key)
	if err != nil {
		return err
	}

	if len(resp.Kvs) == 0 {
		return fmt.Errorf("failed to find key, %s", key)
	}

	err = json.Unmarshal(resp.Kvs[0].Value, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	return nil
}

func (etcd *EtcdClient) GetAllDataWithPrefix(prefix string) (map[string]string, error) {
	resp, err := etcd.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get data from etcd with prefix %s: %v", prefix, err)
	}

	data := make(map[string]string)
	for _, kv := range resp.Kvs {
		data[string(kv.Key)] = string(kv.Value)
	}
	return data, nil
}

func (etcd *EtcdClient) DeleteData(key string) error {
	_, err := etcd.Delete(context.Background(), key)
	if err != nil {
		return fmt.Errorf("failed to delete data from etcd: %v", err)
	}
	return nil
}

func (etcd *EtcdClient) Heartbeat(key string, value interface{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := etcd.PutData(key, value)
		if err != nil {
			log.Printf("Failed to send heartbeat for key %s: %v", key, err)
		} else {
			log.Printf("Heartbeat sent for key %s", key)
		}
	}
}

func (etcd *EtcdClient) Close() {
	etcd.Client.Close()
}
