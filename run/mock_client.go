package main

import (
	"ZTWssProxy/configs"
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := fmt.Sprintf("https://%s", configs.ClientConnAddr)
	data := []byte(`{"key":"value"}`)

	// 创建 HTTP 客户端
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

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
