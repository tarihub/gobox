package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func uploadFile(url string, filePath string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建请求
	request, err := http.NewRequest("POST", url, file)
	if err != nil {
		return err
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/octet-stream")

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 读取响应体
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// 输出响应状态码和响应体
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Response Body: %s\n", string(responseBody))

	return nil
}

func main() {
	// 从命令行参数获取 URL 和要上传的文件路径
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <url> <file>", os.Args[0])
		os.Exit(1)
	}
	url := os.Args[1]
	filePath := os.Args[2]

	err := uploadFile(url, filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("File uploaded successfully")
}
