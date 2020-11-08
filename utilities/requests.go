package utilities

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"time"
)

// requestInfo aggregates pointers affected by functional options for convenience and readability.
type requestInfo struct {
	request  *fasthttp.Request
	response *fasthttp.Response
	timeout  time.Duration
}

// RequestGET provides a basic foundation for performing HTTP GET requests, with optional functionalities provided by
// variadic functional arguments. Note that returned non-errored responses should (in best practice) be manually
// released when finished via fasthttp.
func RequestGET(client *fasthttp.Client, uri string, opts ...func(*requestInfo)) (*fasthttp.Response, error) {
	timeout := time.Minute // default timeout (since 0 timeout breaks)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// apply functional options to request
	info := requestInfo{req, resp, timeout}
	for _, opt := range opts {
		opt(&info)
	}

	// perform actual HTTP request
	req.Header.SetMethod("GET")
	req.SetRequestURI(uri)
	if err := client.DoTimeout(req, resp, timeout); err != nil { // wrap error and return
		return nil, errors.Wrap(err, "HTTP GET error")
	}

	return resp, nil
}

// RequestPOST provides a basic foundation for performing HTTP POST requests, with optional functionalities provided by
// variadic functional arguments. Note that returned non-errored responses should (in best practice) be manually
// released when finished via fasthttp.
func RequestPOST(client *fasthttp.Client, uri string, body string, opts ...func(*requestInfo)) (*fasthttp.Response, error) {
	timeout := time.Minute // default timeout (since 0 timeout breaks)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// apply functional options to request
	info := requestInfo{req, resp, timeout}
	for _, opt := range opts {
		opt(&info)
	}

	// perform actual HTTP request
	req.Header.SetMethod("POST")
	req.SetBodyString(body)
	req.SetRequestURI(uri)
	if err := client.DoTimeout(req, resp, timeout); err != nil { // wrap error and return
		return nil, errors.Wrap(err, "HTTP POST error")
	}

	return resp, nil
}
