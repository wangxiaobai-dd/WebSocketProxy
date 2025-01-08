package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"ZTWssProxy/configs"
	"ZTWssProxy/internal/proxyserver"
)

func main() {
	t := proxyserver.Token{
		LoginTempID: configs.TestLoginTempID,
		AccID:       67890,
		ZoneID:      1,
		GateIp:      "127.0.0.1",
		GatePort:    configs.TestGatePort,
	}

	data, err := json.Marshal(t)
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	url := fmt.Sprintf("http://%s/token", configs.GameTokenAddr)
	client := &http.Client{}

	// 发送 POST 请求
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应状态码和内容
	fmt.Println("Response Status Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(body))
}
