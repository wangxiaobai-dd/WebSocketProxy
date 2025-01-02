package proxyserver

const (
	proxyAddr     = "localhost:8080" // 代理服务器监听地址
	recvTokenAddr = "localhost:8081" // 游戏服务器向代理服务器发送 token 的 HTTP 地址
)

type Server struct {
}
