package httpclient

import (
	"time"

	"log/slog"

	"github.com/go-resty/resty/v2"
)

// HttpClient defines the contract for making HTTP requests
type HttpClient interface {
	Get(url string, headers map[string]string, timeout ...time.Duration) (*resty.Response, error)
	Post(url string, body interface{}, headers map[string]string, timeout ...time.Duration) (*resty.Response, error)
	Put(url string, body interface{}, headers map[string]string, timeout ...time.Duration) (*resty.Response, error)
	Delete(url string, headers map[string]string, timeout ...time.Duration) (*resty.Response, error)
}

// RestClient implements HttpClient using Resty
type RestClient struct {
	client         *resty.Client
	logger         *slog.Logger
	defaultTimeout time.Duration
}

// NewRestClient initializes a Resty client with logging and a default timeout
func NewRestClient(defaultTimeout time.Duration, logger *slog.Logger) *RestClient {
	client := resty.New().
		SetTimeout(defaultTimeout) // Set default timeout
		//SetRetryCount(3).
		//SetRetryWaitTime(2 * time.Second).
		//SetRetryMaxWaitTime(10 * time.Second)

	rc := &RestClient{
		client:         client,
		logger:         logger,
		defaultTimeout: defaultTimeout,
	}

	// Enable request/response logging
	client.OnBeforeRequest(rc.logRequest)
	client.OnAfterResponse(rc.logResponse)

	return rc
}

// logRequest logs outgoing HTTP requests
func (r *RestClient) logRequest(c *resty.Client, req *resty.Request) error {
	r.logger.Info("HTTP Request",
		slog.String("method", req.Method),
		slog.String("url", req.URL),
		slog.Any("headers", req.Header),
	)
	return nil
}

// logResponse logs incoming HTTP responses
func (r *RestClient) logResponse(c *resty.Client, resp *resty.Response) error {
	if resp.IsError() {
		r.logger.Error("HTTP Response Error",
			slog.Int("status", resp.StatusCode()),
			slog.String("url", resp.Request.URL),
			slog.String("body", resp.String()),
		)
	} else {
		r.logger.Info("HTTP Response",
			slog.Int("status", resp.StatusCode()),
			slog.String("url", resp.Request.URL),
			slog.String("body", resp.String()),
		)
	}
	return nil
}

// executeRequest runs the request with an optional timeout
func (r *RestClient) executeRequest(method, url string, body interface{}, headers map[string]string, timeout ...time.Duration) (*resty.Response, error) {
	req := r.client.R()

	// Apply headers if provided
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}

	// Apply request body if applicable
	if body != nil {
		req.SetBody(body)
	}

	// Clone the client and apply the timeout dynamically
	client := r.client
	if len(timeout) > 0 {
		client = resty.New().SetTimeout(timeout[0])
	}

	// Execute request based on the method
	switch method {
	case "GET":
		return client.R().SetHeaders(headers).Get(url)
	case "POST":
		return client.R().SetHeaders(headers).SetBody(body).Post(url)
	case "PUT":
		return client.R().SetHeaders(headers).SetBody(body).Put(url)
	case "DELETE":
		return client.R().SetHeaders(headers).Delete(url)
	default:
		return nil, nil
	}
}

// Get makes a GET request with optional headers and timeout
func (r *RestClient) Get(url string, headers map[string]string, timeout ...time.Duration) (*resty.Response, error) {
	return r.executeRequest("GET", url, nil, headers, timeout...)
}

// Post makes a POST request with optional headers and timeout
func (r *RestClient) Post(url string, body interface{}, headers map[string]string, timeout ...time.Duration) (*resty.Response, error) {
	return r.executeRequest("POST", url, body, headers, timeout...)
}

// Put makes a PUT request with optional headers and timeout
func (r *RestClient) Put(url string, body interface{}, headers map[string]string, timeout ...time.Duration) (*resty.Response, error) {
	return r.executeRequest("PUT", url, body, headers, timeout...)
}

// Delete makes a DELETE request with optional headers and timeout
func (r *RestClient) Delete(url string, headers map[string]string, timeout ...time.Duration) (*resty.Response, error) {
	return r.executeRequest("DELETE", url, nil, headers, timeout...)
}
