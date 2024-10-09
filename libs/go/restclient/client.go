package restclient

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type restclient struct {
	client               *resty.Client
	baseUrl              string
	CircuitBreakerConfig CircuitBreakerConfig
	// MetricsConfig        MetricsConfig
}

func New(cfg Config) Client {
	client := resty.New()
	client.SetBaseURL(cfg.BaseUrl)
	client.SetRetryCount(cfg.Retries)
	client.SetDebug(cfg.DebugMode)

	if cfg.CustomTransport != nil {
		client.SetTransport(cfg.CustomTransport)
	}

	if cfg.TimeoutMs != nil {
		client.SetTimeout(time.Duration(*cfg.TimeoutMs) * time.Millisecond)
	}

	if cfg.RetryWaitTimeMs != nil {
		retryWaitTime := *cfg.RetryWaitTimeMs
		if retryWaitTime <= 0 {
			retryWaitTime = 1
		}
		client.SetRetryWaitTime(time.Duration(retryWaitTime) * time.Millisecond)
	}

	client.OnBeforeRequest(onBeforeRequestHook)

	/** Circuit Breaker */
	var circuitBreakerCfg CircuitBreakerConfig
	if cfg.CircuitBreakerConfig != nil && cfg.CircuitBreakerConfig.Enabled {
		circuitBreakerCfg.Enabled = true
		circuitBreakerCfg.Timeout = cfg.CircuitBreakerConfig.Timeout
		circuitBreakerCfg.FailedRequests = cfg.CircuitBreakerConfig.FailedRequests
		circuitBreakerCfg.FailureRatioAllowed = cfg.CircuitBreakerConfig.FailureRatioAllowed
		circuitBreakerCfg.MaxRequests = cfg.CircuitBreakerConfig.MaxRequests
	}

	/** Datadog Metrics */
	// var metricsCfg MetricsConfig
	// if cfg.MetricsConfig != nil && cfg.MetricsConfig.Enabled {
	// 	metricsCfg.Enabled = true
	// }

	rc := restclient{
		client:               client,
		CircuitBreakerConfig: circuitBreakerCfg,
		baseUrl:              cfg.BaseUrl,
		// MetricsConfig:        metricsCfg,
	}

	return &rc
}

func (rc *restclient) GET(urlFormat string, opts ...EndpointOption) Endpoint {

	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodGet,
		headers:   make(http.Header),
		opts:      opts,
	}

	if rc.CircuitBreakerConfig.Enabled {
		endpoint = EndpointWithCircuitBreaker(rc.CircuitBreakerConfig, endpoint)
	}

	// if rc.MetricsConfig.Enabled {
	// 	requestUrl := rc.baseUrl + urlFormat
	// 	endpoint = EndpointWithMetrics(requestUrl, endpoint)
	// }

	return endpoint
}

func (rc *restclient) POST(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodPost,
		headers:   make(http.Header),
		opts:      opts,
	}

	if rc.CircuitBreakerConfig.Enabled {
		endpoint = EndpointWithCircuitBreaker(rc.CircuitBreakerConfig, endpoint)
	}

	// if rc.MetricsConfig.Enabled {
	// 	requestUrl := rc.baseUrl + urlFormat
	// 	endpoint = EndpointWithMetrics(requestUrl, endpoint)
	// }

	return endpoint
}

func (rc *restclient) PUT(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodPut,
		headers:   make(http.Header),
		opts:      opts,
	}

	if rc.CircuitBreakerConfig.Enabled {
		endpoint = EndpointWithCircuitBreaker(rc.CircuitBreakerConfig, endpoint)
	}

	// if rc.MetricsConfig.Enabled {
	// 	requestUrl := rc.baseUrl + urlFormat
	// 	endpoint = EndpointWithMetrics(requestUrl, endpoint)
	// }

	return endpoint
}

func (rc *restclient) PATCH(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodPatch,
		headers:   make(http.Header),
		opts:      opts,
	}

	if rc.CircuitBreakerConfig.Enabled {
		endpoint = EndpointWithCircuitBreaker(rc.CircuitBreakerConfig, endpoint)
	}

	// if rc.MetricsConfig.Enabled {
	// 	requestUrl := rc.baseUrl + urlFormat
	// 	endpoint = EndpointWithMetrics(requestUrl, endpoint)
	// }

	return endpoint
}

func (rc *restclient) DELETE(urlFormat string, opts ...EndpointOption) Endpoint {
	var endpoint Endpoint = &endpoint{
		client:    rc.client,
		urlFormat: urlFormat,
		method:    http.MethodDelete,
		headers:   make(http.Header),
		opts:      opts,
	}

	if rc.CircuitBreakerConfig.Enabled {
		endpoint = EndpointWithCircuitBreaker(rc.CircuitBreakerConfig, endpoint)
	}

	// if rc.MetricsConfig.Enabled {
	// 	requestUrl := rc.baseUrl + urlFormat
	// 	endpoint = EndpointWithMetrics(requestUrl, endpoint)
	// }

	return endpoint
}

func onBeforeRequestHook(c *resty.Client, r *resty.Request) error {
	return nil
}
