package hotelbyte

import (
	"log"

	"github.com/bytedance/sonic"
)

func ToJSON(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return "" // 兼容 nil 值，不要序列化成 null
	}
	b, err := sonic.Marshal(v)
	if err != nil {
		log.Printf("ToJSONString failed(%v) from value(%v)\n", err, v)
		return ""
	}
	return string(b)
}
