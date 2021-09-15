package utils

import (
	"io"
	"net/http"
)

// 微信 https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Receiving_standard_messages.html
//     https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Receiving_event_pushes.html
// 企业微信 https://work.weixin.qq.com/api/doc/90000/90135/90238
// 跳过相关校验之后的处理回调
type XmlHandlerFunc func(http.ResponseWriter, *http.Request, []byte) error

func HttpAbort(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	io.WriteString(w, http.StatusText(code))
}

func HttpAbortBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}
