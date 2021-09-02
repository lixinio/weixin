package utils

import (
	"errors"
	"fmt"
)

var (
	ErrorAccessToken = errors.New("access token error")
	ErrorSystemBusy  = errors.New("system busy")
	ErrorWeixinError = errors.New("system busy")
)

type WeixinErrorInterface interface {
	WeixinErrorCode() int64
	WeixinErrorMessage() string
	GetWeixinError() error
}

type WeixinError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// @WeixinErrorInterface
func (we *WeixinError) WeixinErrorCode() int64 {
	return we.ErrCode
}

// @WeixinErrorInterface
func (we *WeixinError) WeixinErrorMessage() string {
	return we.ErrMsg
}

// @WeixinErrorInterface
func (we *WeixinError) GetWeixinError() error {
	return we
}

// @error
func (we *WeixinError) Error() string {
	return fmt.Sprintf("%d: %s", we.ErrCode, we.ErrMsg)
}
