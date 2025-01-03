package main

import (
	"log"
	"os"
	"time"

	"ZTWssProxy/internal/proxyserver"
)

func initLog() {
	currentTime := time.Now()
	logFileName := currentTime.Format("2006-01-02") + ".log" // 使用日期作为文件名

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// 设置日志输出到文件
	log.SetOutput(logFile)

	// 设置日志的前缀和格式
	log.SetPrefix("[INFO] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {

}

func main() {
	server := proxyserver.NewProxyServer()
	server.Run()
}
