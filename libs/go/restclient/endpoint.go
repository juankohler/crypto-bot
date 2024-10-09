package restclient

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type endpoint struct {
	urlFormat string
	method    string
	headers   http.Header
	client    *resty.Client
	opts      []EndpointOption
}

type EndpointOption func(p *request)

func (e *endpoint) DoRequest(ctx context.Context, opts ...EndpointOption) Response {
	return newRequest(e).Do(ctx, opts...)
}

func (e *endpoint) Request() Request {
	return newRequest(e)
}

func UrlParam(key string, value interface{}) EndpointOption {
	return func(p *request) {
		p.URL = strings.ReplaceAll(p.URL, "{"+key+"}", fmt.Sprint(value))
	}
}

func QueryParam(key, value string) EndpointOption {
	return func(p *request) {
		p.SetQueryParam(key, value)
	}
}

func QueryParamList(key string, values []string) EndpointOption {
	return func(p *request) {
		p.SetQueryParam(key, strings.Join(values, ","))
	}
}

func SetFormData(data map[string]string) EndpointOption {
	return func(p *request) {
		p.SetFormData(data)
	}
}

func SetBasicAuth(username string, password string) EndpointOption {
	return func(p *request) {
		p.SetBasicAuth(username, password)
	}
}

func Body(body interface{}) EndpointOption {
	return func(p *request) {
		p.SetBody(body)
	}
}

func Header(key string, value string) EndpointOption {
	return func(p *request) {
		p.SetHeader(key, value)
	}
}
func Headers(hs map[string]string) EndpointOption {
	return func(p *request) {
		p.SetHeaders(hs)
	}
}

func BasicAuth(username, password string) EndpointOption {
	return func(p *request) {
		p.SetBasicAuth(username, password)
	}
}

func FailAt(failAt FailAtFunc) EndpointOption {
	return func(p *request) {
		p.SetFailAt(failAt)
	}
}
