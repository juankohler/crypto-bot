package restclient

import (
	"context"
	"net/http"
	"time"

	"github.com/juankohler/crypto-bot/libs/go/errors"
	"github.com/sony/gobreaker"
)

var (
	ErrTypeCircuitBreaker = errors.Define("circuit_breaker")
)

// Config
type CircuitBreakerConfig struct {
	/** True to initialize a circuit breaker */
	Enabled bool
	// Sleep time to consider new requests after FailedRequests or
	// FailureRatioAllowed are reached.
	Timeout time.Duration
	// Subsequent requests with errors allowed before passing to open state.
	FailedRequests uint32
	// Ratio of failures / requests allowed before passing to open state.
	FailureRatioAllowed float64
	// MaxRequests is the maximum number of requests allowed to pass through
	// when the CircuitBreaker is half-open.
	// If MaxRequests is 0, the CircuitBreaker allows only 1 request.
	MaxRequests uint32
}

// Wrapper for endpoints with circuit breaker
type circuitBreakerEndpoint struct {
	endpoint       Endpoint
	circuitBreaker *gobreaker.CircuitBreaker
}

func EndpointWithCircuitBreaker(
	config CircuitBreakerConfig,
	endpoint Endpoint,
) Endpoint {
	return &circuitBreakerEndpoint{
		endpoint: endpoint,
		circuitBreaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Timeout:     config.Timeout,
			MaxRequests: config.MaxRequests,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= config.FailedRequests && failureRatio >= config.FailureRatioAllowed
			},
		}),
	}
}

func (e *circuitBreakerEndpoint) DoRequest(ctx context.Context, opts ...EndpointOption) Response {
	res, err := e.circuitBreaker.Execute(func() (interface{}, error) {
		// If circuit breaker is in open or half-open state, request will be
		// executed here, so we'll get the real response.
		// In case circuit breaker is in closed state, this callback will never
		// be executed, returning a Circuit Breaker error.
		// So it's safe to return the real response without an error.
		res := e.endpoint.DoRequest(ctx, opts...)

		return res, res.Err()
	})

	if res, ok := res.(Response); ok {
		if err != nil {
			return &circuitBreakerErrorResponse{res, err}
		}

		return res
	}

	return &circuitBreakerErrorResponse{err: errors.Wrap(ErrTypeCircuitBreaker, err, "circuit breaker blocked the request")}
}

func (e *circuitBreakerEndpoint) Request() Request {
	return e.endpoint.Request()
}

// Response with circuit breaker error
type circuitBreakerErrorResponse struct {
	res Response
	err error
}

func (r *circuitBreakerErrorResponse) Body() []byte {
	if r.res != nil {
		return r.res.Body()
	}
	return []byte{}
}

func (r *circuitBreakerErrorResponse) Status() string {
	if r.res != nil {
		return r.res.Status()
	}
	return ""
}

func (r *circuitBreakerErrorResponse) StatusCode() int {
	if r.res != nil {
		return r.res.StatusCode()
	}
	return -1
}

func (r *circuitBreakerErrorResponse) Header() http.Header {
	if r.res != nil {
		return r.res.Header()
	}
	return http.Header{}
}

func (r *circuitBreakerErrorResponse) Err() error {
	return r.err
}
