package types

import (
	"fmt"
	"github.com/bytedance/sonic"
	"net/http"
	"net/url"
	"reflect"
)

type Response[T any] struct {
	BizError
	Data   *T          `json:"data"`
	Header http.Header `json:"header"`
}

func NewResponseData[T any](r *HttpResponse) (*T, error) {
	v, err := NewResponse[T](r)
	if err != nil {
		return nil, err
	}

	// 如果 T 类型中有 header 字段，则将 v.Header赋值给 T
	if v.Data != nil {
		setHeaderField(v.Data, v.Header)
	}

	return v.Data, nil
}

// 使用反射为结构体设置header字段
func setHeaderField[T any](data *T, header http.Header) {
	if data == nil {
		return
	}

	// 获取值的反射对象
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	// 遍历结构体字段
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// 检查字段名是否为 "header"
		if fieldName == "header" {
			// 检查字段类型是否是 http.Header (map[string][]string)
			if field.Type.Kind() == reflect.Map &&
			   field.Type.Key().Kind() == reflect.String &&
			   field.Type.Elem().Kind() == reflect.Slice &&
			   field.Type.Elem().Elem().Kind() == reflect.String {

				headerValue := val.Field(i)
				if headerValue.IsValid() && headerValue.CanSet() {
					headerValue.Set(reflect.ValueOf(header))
				}
				break
			}
		}
	}
}

func NewResponse[T any](r *HttpResponse) (*Response[T], error) {
	if len(r.Body) == 0 {
		return nil, NewBizErr(int32(r.StatusCode), "Service Unavailable")
	}
	var response Response[T]
	err := sonic.Unmarshal(r.Body, &response)
	if err != nil {
		return nil, fmt.Errorf("invalid body %w", err)
	}
	if response.Code != 0 {
		return nil, &response.BizError
	}
	response.Header = r.Headers
	return &response, nil
}

// HttpRequest represents HTTP request
type HttpRequest struct {
	Method  string
	Path    string
	Query   url.Values
	Headers map[string]string
	Body    interface{}
}

// HttpResponse represents HTTP response
type HttpResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}
