package types

import (
	"fmt"
	"github.com/bytedance/sonic"
)

type Response[T any] struct {
	BizError
	Data *T `json:"data"`
}

func NewResponse[T any](status int, body []byte) (*T, error) {
    if len(body) == 0 {
        return nil, NewBizErr(int32(status), "Service Unavailable")
    }
    var response Response[T]
    err := sonic.Unmarshal(body, &response)
    if err != nil {
        return nil, fmt.Errorf("invalid body %w", err)
    }
    if response.Code != 0 {
        return nil, &response.BizError
    }
    // 防御性检查：确保 data 字段存在，避免后续解引用导致的 panic
    if response.Data == nil {
        return nil, NewBizErr(int32(status), "Invalid Response: missing data")
    }
    return response.Data, nil
}
