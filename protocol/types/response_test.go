package types

import (
	"net/http"
	"testing"
)

// 定义一个包含header字段的结构体来测试
type TestResponseWithHeader struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Header http.Header `json:"header"`
}

// 定义一个不包含header字段的结构体来测试
type TestResponseWithoutHeader struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestNewResponseDataWithHeaderField(t *testing.T) {
	// 模拟HTTP响应
	responseBody := `{"code":0,"data":{"id":123,"name":"test","header":null}}`
	response := &HttpResponse{
		StatusCode: 200,
		Headers:    http.Header{
			"X-Custom-Header": []string{"value1"},
			"X-Request-ID":    []string{"req-123"},
		},
		Body: []byte(responseBody),
	}

	result, err := NewResponseData[TestResponseWithHeader](response)
	if err != nil {
		t.Fatalf("NewResponseData failed: %v", err)
	}

	// 验证基本数据
	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}
	if result.Name != "test" {
		t.Errorf("Expected Name 'test', got %s", result.Name)
	}

	// 验证header字段是否被正确赋值
	if result.Header == nil {
		t.Error("Expected header to be set, but got nil")
	}
	if result.Header.Get("X-Custom-Header") != "value1" {
		t.Errorf("Expected X-Custom-Header to be 'value1', got '%s'", result.Header.Get("X-Custom-Header"))
	}
	if result.Header.Get("X-Request-ID") != "req-123" {
		t.Errorf("Expected X-Request-ID to be 'req-123', got '%s'", result.Header.Get("X-Request-ID"))
	}
}

func TestNewResponseDataWithoutHeaderField(t *testing.T) {
	// 模拟HTTP响应
	responseBody := `{"code":0,"data":{"id":456,"name":"test-no-header"}}`
	response := &HttpResponse{
		StatusCode: 200,
		Headers:    http.Header{
			"X-Custom-Header": []string{"value2"},
		},
		Body: []byte(responseBody),
	}

	result, err := NewResponseData[TestResponseWithoutHeader](response)
	if err != nil {
		t.Fatalf("NewResponseData failed: %v", err)
	}

	// 验证基本数据
	if result.ID != 456 {
		t.Errorf("Expected ID 456, got %d", result.ID)
	}
	if result.Name != "test-no-header" {
		t.Errorf("Expected Name 'test-no-header', got %s", result.Name)
	}

	// 验证不包含header字段的结构体不会出错
	// 应该正常运行，不会尝试设置header字段
}

func TestNewResponseDataWithErrorResponse(t *testing.T) {
	// 模拟错误响应
	responseBody := `{"code":500,"msg":"Internal Server Error","data":null}`
	response := &HttpResponse{
		StatusCode: 500,
		Headers:    http.Header{},
		Body:       []byte(responseBody),
	}

	result, err := NewResponseData[TestResponseWithHeader](response)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if result != nil {
		t.Error("Expected nil result on error, but got non-nil")
	}
}

func TestNewResponseDataWithEmptyBody(t *testing.T) {
	response := &HttpResponse{
		StatusCode: 204,
		Headers:    http.Header{},
		Body:       []byte(""),
	}

	result, err := NewResponseData[TestResponseWithHeader](response)
	if err == nil {
		t.Error("Expected error for empty body, but got nil")
	}
	if result != nil {
		t.Error("Expected nil result on error, but got non-nil")
	}
}

func TestSetHeaderField(t *testing.T) {
	// 测试设置header字段的辅助函数
	header := http.Header{
		"Content-Type": []string{"application/json"},
		"Authorization": []string{"Bearer token"},
	}

	data := &TestResponseWithHeader{}
	setHeaderField(data, header)

	if data.Header == nil {
		t.Error("Expected header to be set, but got nil")
	}
	if data.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', got '%s'", data.Header.Get("Content-Type"))
	}
	if data.Header.Get("Authorization") != "Bearer token" {
		t.Errorf("Expected Authorization to be 'Bearer token', got '%s'", data.Header.Get("Authorization"))
	}
}