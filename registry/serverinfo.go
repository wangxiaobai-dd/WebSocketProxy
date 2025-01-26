package registry

import "fmt"

type ServerInfo struct {
	ServerID     int    `json:"serverID"`
	ServerIP     string `json:"serverIP"`
	ServerDomain string `json:"serverDomain"`
	TokenPort    int    `json:"tokenPort"`
	ClientPort   int    `json:"clientPort"`
	ConnNum      int    `json:"connNum"`
	SecureFlag   bool   `json:"secureFlag"`
}

func getServerConnSetKey(prefix string) string {
	return fmt.Sprintf("connset-%s", prefix)
}

func getServerKey(prefix string, serverID int) string {
	return fmt.Sprintf("%s:%d", prefix, serverID)
}
