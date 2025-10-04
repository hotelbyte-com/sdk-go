package types

import (
	"errors"
	"fmt"
	"reflect"
)

type BizError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func (e *BizError) String() string {
	return e.Error()
}

func (e *BizError) Error() string {
	if e == nil {
		return "(*BizError)<nil>"
	}
	return fmt.Sprintf("ErrorCode=[%v] Message=[%v]", e.Code, e.Msg)
}

// WithMessage 修改 StatusMessage 并会返回一个新的 BizError
func (e *BizError) WithMessage(msg string) *BizError {
	bizErr := copyErr(e)
	bizErr.Msg = msg
	return bizErr
}

// WithMessagef 修改 StatusMessage 并会返回一个新的 BizError
func (e *BizError) WithMessagef(format string, a ...any) *BizError {
	bizErr := copyErr(e)
	bizErr.Msg = fmt.Sprintf(format, a...)
	return bizErr
}
func (e *BizError) StatusCode() int32 {
	if e == nil {
		return 0
	}
	return e.Code
}

func (e *BizError) StatusMessage() string {
	return e.Msg
}

// Is implement errors.Is
func (e *BizError) Is(oe error) bool {
	return e.HasSameCode(oe)
}

func (e *BizError) HasSameCode(err error) (ok bool) {
	if errors.Is(err, e) {
		return true
	}
	var bizErr *BizError

	if bizErr, ok = UnwrapForBizErr(err); ok {
		return e.Code == bizErr.Code
	}
	return false
}

func New(statusCode int32, msg string) *BizError {
	return &BizError{Code: statusCode, Msg: msg}
}

func copyErr(e *BizError) *BizError {
	return &BizError{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

func CastBizErr(err error) (bizErr *BizError, ok bool) {
	if IsNil(err) {
		return nil, false
	}
	ok = errors.As(err, &bizErr)
	return bizErr, ok
}

func UnwrapForBizErr(err error) (*BizError, bool) {
	var bizErrType *BizError
	if !errors.As(err, &bizErrType) {
		return nil, false
	}
	maxWrap := 100
	for i := 1; !IsNil(err) && i <= maxWrap; i++ {
		var bizErr *BizError
		if errors.As(err, &bizErr) {
			return bizErr, true
		}
		err = errors.Unwrap(err)
	}
	return nil, false
}

func IsNil(err error) bool {
	if err == nil {
		return true
	}
	t := reflect.ValueOf(err)
	return t.Kind() == reflect.Ptr && t.IsNil()
}
