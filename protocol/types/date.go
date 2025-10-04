package types

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"

	"github.com/spf13/cast"
)

// DateInt compatible with string & int64, e.g. 20230101, 2023-01-01
type DateInt int

func (d DateInt) Int64() int64 {
	return int64(d)
}
func (d DateInt) Format(layout string) string {
	val := int(d)
	if layout == "" {
		return strconv.Itoa(val)
	}
	if val <= 0 || val < 10000101 {
		return ""
	}
	year := val / 10000
	month := (val % 10000) / 100
	day := val % 100
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return ""
	}
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return ""
	}
	return t.Format(layout)
}

func (d DateInt) Sub(in DateInt) int {
	return int(d.ToTime().Sub(in.ToTime()).Hours() / 24)
}

func (d DateInt) AddDays(days int) DateInt {
	v := d.ToTime().AddDate(0, 0, days)
	return NewDateIntFromTime(v)
}
func (d DateInt) SubDays(days int) DateInt {
	v := d.ToTime().AddDate(0, 0, -days)
	return NewDateIntFromTime(v)
}
func (d DateInt) AddMonths(months int) DateInt {
	v := d.ToTime().AddDate(0, months, 0)
	return NewDateIntFromTime(v)
}
func (d DateInt) SubMonths(months int) DateInt {
	v := d.ToTime().AddDate(0, -months, 0)
	return NewDateIntFromTime(v)
}
func (d DateInt) ToTime() time.Time {
	return time.Date(d.Year(), time.Month(d.Month()), d.Day(), 0, 0, 0, 0, time.UTC)
}
func (d DateInt) Year() int {
	return int(int64(d) / 10000)
}
func (d DateInt) Month() int {
	return int((int64(d) % 10000) / 100)
}
func (d DateInt) Day() int {
	return int(int64(d) % 100)
}

func compatibleParseDateOnlyStr(hint, s string) time.Time {
	if hint != "" {
		if t, err := time.Parse(hint, s); err == nil {
			return t
		}
	}
	if t, err := time.Parse("20060102", s); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t
	}
	return time.Time{}
}

func ParseDateInt(ctx context.Context, date string) DateInt {
	t := compatibleParseDateOnlyStr("", date)
	if t.IsZero() {
		return 0
	}
	return DateInt(cast.ToInt(t.Format("20060102")))
}
func ParseDateIntWithHint(ctx context.Context, date, hint string) DateInt {
	t := compatibleParseDateOnlyStr(hint, date)
	if t.IsZero() {
		return 0
	}
	return DateInt(cast.ToInt(t.Format("20060102")))
}
func NewDateIntFromTime(t time.Time) DateInt {
	if t.IsZero() {
		return 0
	}
	return DateInt(cast.ToInt(t.Year()*10000 + int(t.Month())*100 + t.Day()))
}
func NewDateIntFromDateString(s string) DateInt {
	vs := strings.Split(s, "-")
	if len(vs) == 3 {
		return DateInt(cast.ToInt64(vs[0])*10000 + cast.ToInt64(vs[1])*100 + cast.ToInt64(vs[2]))
	}
	return NewDateIntFromTime(compatibleParseDateOnlyStr("20060102", s))
}

// parseBuiltinFunction 解析内置函数，支持 now(), NOW() + N, now() - N 等格式
func parseBuiltinFunction(str string) (time.Time, bool) {
	str = strings.TrimSpace(strings.ToLower(str))
	if !strings.HasPrefix(str, "now()") {
		return time.Time{}, false
	}
	now := time.Now()
	if str == "now()" {
		return now, true
	}
	remaining := strings.TrimSpace(str[5:])
	if remaining == "" {
		return now, true
	}
	var op rune
	var numStr string
	for i, r := range remaining {
		if r == '+' || r == '-' {
			op = r
			numStr = strings.TrimSpace(remaining[i+1:])
			break
		}
	}
	if op == 0 || numStr == "" {
		return time.Time{}, false
	}
	offset, err := strconv.Atoi(numStr)
	if err != nil {
		return time.Time{}, false
	}
	if op == '+' {
		return now.AddDate(0, 0, offset), true
	}
	return now.AddDate(0, 0, -offset), true
}

func (d *DateInt) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		str := string(data[1 : len(data)-1])
		if str == "" {
			*d = 0
			return nil
		}
		if strings.Contains(str, "(") {
			if t, ok := parseBuiltinFunction(str); ok {
				*d = DateInt(cast.ToInt(t.Format("20060102")))
				return nil
			}
		}
		switch len(str) {
		case 10:
			if t, err := time.Parse("2006-01-02", str); err == nil {
				*d = NewDateIntFromTime(t)
				return nil
			}
		case 8:
			if _, err := strconv.Atoi(str); err == nil {
				if t, err := time.Parse("20060102", str); err == nil {
					*d = NewDateIntFromTime(t)
					return nil
				}
			}
		}
		if t := compatibleParseDateOnlyStr("", str); !t.IsZero() {
			*d = NewDateIntFromTime(t)
			return nil
		}
		return fmt.Errorf("invalid date string: %s", str)
	}
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*d = DateInt(num)
		return nil
	}
	return fmt.Errorf("invalid date format: %s", string(data))
}

// MarshalJSON YYYY-MM-DD, "2006-01-02"
func (d DateInt) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(d.Format("2006-01-02"))
}
