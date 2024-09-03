package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PostToServer(data interface{}) (map[string]interface{}, error) {
	// 将数据序列化为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON 序列化错误: %v", err)
	}

	// 创建一个 POST 请求，将 JSON 数据作为请求体发送
	resp, err := http.Post("http://localhost:8080/process", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("发送请求错误: %v", err)
	}
	defer resp.Body.Close()

	// 读取服务器响应
	body, err := io.ReadAll(resp.Body) // 使用 io.ReadAll 替代 ioutil.ReadAll
	if err != nil {
		return nil, fmt.Errorf("读取响应错误: %v", err)
	}

	// 将服务器响应解析为 map[string]interface{}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("JSON 解析错误: %v", err)
	}

	return result, nil
}
