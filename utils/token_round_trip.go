package utils

import (
	"context"
	"fmt"
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

// 在Trace的时候， 移除access-token / secret
// 	secret : https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html

type (
	AccessTokenStripTransport struct {
		defaultKey string
		Base       http.RoundTripper
	}
	stripKeyContext int
)

var stripKeyContextKey = stripKeyContext(0)

func NewStripContext(ctx context.Context, keys ...string) context.Context {
	return context.WithValue(ctx, stripKeyContextKey, keys)
}

func parseStripContext(ctx context.Context) ([]string, bool) {
	v := ctx.Value(stripKeyContextKey)
	if v != nil {
		if k, ok := v.([]string); ok && len(k) > 0 {
			return k, true
		}
	}

	return nil, false
}

func (t *AccessTokenStripTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Base.RoundTrip(req)
	ctx := req.Context()

	span := trace.FromContext(ctx)
	if span == nil {
		return resp, err
	}

	u := req.URL
	q := u.Query()

	// 如果存在， 重置
	edit := false
	if q.Get(t.defaultKey) != "" {
		q.Set(t.defaultKey, "")
		edit = true
	}

	if stripKeys, ok := parseStripContext(ctx); ok {
		for _, stripKey := range stripKeys {
			if q.Get(stripKey) != "" {
				q.Set(stripKey, "")
				edit = true
			}
		}
	}

	if edit {
		// 覆盖原来的Url
		span.AddAttributes(trace.StringAttribute(
			ochttp.URLAttribute,
			fmt.Sprintf("%s://%s%s?%s", u.Scheme, u.Host, u.Path, q.Encode()),
		))
	}
	return resp, err
}

func NewAccessTokenStripTransport(defaultKey string) *AccessTokenStripTransport {
	return &AccessTokenStripTransport{
		defaultKey: defaultKey,
		Base:       http.DefaultTransport,
	}
}
