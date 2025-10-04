package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

// ID compatible with string & int64
type ID int64

func NewID(id any) ID {
	return ID(cast.ToInt64(id))
}

type IDs []ID

func (ids IDs) Join(sep string) string {
	return joinInt64Slice(ids.Int64s(), sep)
}

func (ids IDs) Int64s() (out []int64) {
	for _, id := range ids {
		out = append(out, id.Int64())
	}
	return out
}

func (ids *IDs) MarshalJSON() ([]byte, error) {
	if ids == nil || *ids == nil {
		return []byte("null"), nil
	}
	if len(*ids) == 0 {
		return []byte("[]"), nil
	}
	return sonic.Marshal(*ids)
}

func (ids *IDs) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	// 去除前后空白，统一用 b 进行分支判断
	b := bytes.TrimSpace(data)
	if len(b) == 0 {
		return nil
	}
	switch b[0] {
	case '[':
		var tmp []ID
		if err := json.Unmarshal(b, &tmp); err != nil {
			return err
		}
		*ids = tmp
		return nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+', '.', 'e', 'n':
		var tmp int64
		if err := sonic.Unmarshal(b, &tmp); err != nil {
			return err
		}
		if tmp == 0 {
			return nil
		}
		*ids = IDs{ID(tmp)}
		return nil
	case '"':
		var tmp string
		if err := json.Unmarshal(b, &tmp); err != nil {
			return err
		}
		*ids = IDs{}
		for _, id := range strings.Split(tmp, ",") {
			*ids = append(*ids, ID(cast.ToInt64(strings.TrimSpace(id))))
		}
		return nil
	default:
		return fmt.Errorf("unknown id type: %s", b)
	}
}

// UnmarshalJSON for ID supports string and number types
func (id *ID) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		str := string(data[1 : len(data)-1])
		if str == "" {
			*id = 0
			return nil
		}
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID string: %s", str)
		}
		*id = ID(val)
		return nil
	}
	var val int64
	if err := sonic.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("invalid ID number: %s", string(data))
	}
	*id = ID(val)
	return nil
}

// MarshalJSON outputs string to avoid JS big-int issues
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id ID) String() string {
	if id == 0 {
		return ""
	}
	return strconv.FormatInt(int64(id), 10)
}
func (id ID) Int64() int64   { return int64(id) }
func (id ID) Uint64() uint64 { return uint64(id) }
func (id ID) IsZero() bool   { return id == 0 }
func (id ID) Valid() bool    { return id > 0 }

// joinInt64Slice joins int64 slice into string with sep
func joinInt64Slice(nums []int64, sep string) string {
	if len(nums) == 0 {
		return ""
	}
	b := strings.Builder{}
	for i, n := range nums {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(strconv.FormatInt(n, 10))
	}
	return b.String()
}
