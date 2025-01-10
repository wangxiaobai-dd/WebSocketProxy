package configs

const (
	ClientConnAddr  = "localhost:8080" // 客户端连接地址
	GameTokenAddr   = "localhost:8081" // 代理服务器接收游戏服务器发送的 token
	CertFile        = "server_ssl.crt"
	KeyFile         = "server_ssl.key"
	TokenExpireTime = 1000

	TestGateIp      = "localhost:"
	TestGatePort    = 5240
	TestLoginTempID = 12345
)

type ProxyConfig struct {
	ServerID      int      `json:"serverID"`
	ServerIP      string   `json:"serverIP"`
	TokenPort     int      `json:"tokenPort"`
	ClientPort    int      `json:"clientPort"`
	EtcdEndPoints []string `json:"etcdEndPoints"`
}
