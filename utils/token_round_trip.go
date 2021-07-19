package utils

import (
	"fmt"
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

// 在Trace的时候， 移除access-token

type AccessTokenStripTransport struct {
	Base http.RoundTripper
}

func (t *AccessTokenStripTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Base.RoundTrip(req)

	span := trace.FromContext(req.Context())
	if span == nil {
		return nil, err
	}

	u := req.URL
	q := u.Query()
	// 如果存在， 重置
	if q.Get("access_token") != "" {
		q.Set("access_token", "")

		// 覆盖原来的Url
		span.AddAttributes(trace.StringAttribute(
			ochttp.URLAttribute,
			fmt.Sprintf("%s://%s%s?%s", u.Scheme, u.Host, u.Path, q.Encode()),
		))
	}
	return resp, err
}

func newTransport() http.RoundTripper {
	return &ochttp.Transport{
		Base: &AccessTokenStripTransport{
			Base: http.DefaultTransport,
		},
	}
}
