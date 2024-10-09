package restclient

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/juankohler/crypto-bot/libs/go/errors"
)

var (
	ErrTimeout = errors.Define("timeout")
)

type request struct {
	*resty.Request
	endpoint *endpoint
	failAt   FailAtFunc
}

func newRequest(e *endpoint) *request {
	r := request{
		Request:  e.client.R(),
		endpoint: e,
		failAt:   defaultFailAt,
	}

	r.Method = e.method
	r.URL = e.urlFormat
	for _, opt := range e.opts {
		opt(&r)
	}
	return &r
}

func (r *request) Do(ctx context.Context, opts ...EndpointOption) Response {
	for _, opt := range opts {
		opt(r)
	}

	resp := response{}
	resp.Response, resp.err = r.Send()

	if err, ok := resp.err.(net.Error); ok && err.Timeout() {
		resp.err = errors.Wrap(
			ErrTimeout,
			err,
			"http request timeout",
			errors.WithMetadata("error", err),
			errors.WithMetadata("method", strings.ToUpper(r.Method)),
			errors.WithMetadata("url", r.URL),
			errors.WithMetadata("timeout", r.endpoint.client.GetClient().Timeout),
			errors.WithMetadata("attempt", r.Attempt),
		)
	}

	if err := r.failAt(r, &resp); err != nil {
		resp.err = err
	}

	return &resp
}

func (r *request) UrlParam(key string, value interface{}) Request {
	r.URL = strings.ReplaceAll(r.URL, "{"+key+"}", fmt.Sprint(value))
	return r
}

func (r *request) QueryParam(key, value string) Request {
	r.SetQueryParam(key, value)
	return r
}

func (r *request) QueryParamList(key string, values []string) Request {
	r.SetQueryParam(key, strings.Join(values, ","))
	return r
}

func (r *request) Body(body interface{}) Request {
	r.SetBody(body)
	return r
}

func (r *request) Header(key string, value interface{}) Request {
	r.SetHeader(key, fmt.Sprint(value))
	return r

}

func (r *request) Headers(hs map[string]string) Request {
	r.SetHeaders(hs)
	return r
}

func (r *request) BasicAuth(username, password string) Request {
	r.SetBasicAuth(username, password)
	return r
}

func (r *request) SetFailAt(failAt FailAtFunc) Request {
	r.failAt = failAt
	return r
}
