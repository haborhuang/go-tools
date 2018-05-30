package types

import "fmt"

type ErrRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *ErrRes) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("(%d)%s", e.ErrCode, e.ErrMsg)
}

func (e *ErrRes) Err() error {
	if e == nil || e.ErrCode == 0 {
		return nil
	}

	return e
}

func (e *ErrRes) IsExpiredTokenErr() bool {
	return e.Err() != nil && (e.ErrCode == 40014 || e.ErrCode == 42001)
}
