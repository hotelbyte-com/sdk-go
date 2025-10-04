package types

import (
	"fmt"
	"github.com/bytedance/sonic"
)

type Response[T any] struct {
	BizError
	Data *T `json:"data"`
}

func NewResponse[T any](body []byte) (*T, error) {
	var response Response[T]
	err := sonic.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("invalid body %w", err)
	}
	if response.Code != 0 {
		return nil, &response.BizError
	}
	return response.Data, nil
}
