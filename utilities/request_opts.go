package utilities

import "time"

// OptHeader attaches the passed key, value header pair to the request.
func OptHeader(k string, v string) func(*requestInfo) {
	return func(r *requestInfo) {
		r.request.Header.Set(k, v)
	}
}

// OptHeaders attaches a set of key, value header pairs to the request.
func OptHeaders(h map[string]string) func(*requestInfo) {
	return func(r *requestInfo) {
		for k, v := range h {
			r.request.Header.Set(k, v)
		}
	}
}

// OptCookie attaches the passed key, value cookie to the request.
func OptCookie(k string, v string) func(*requestInfo) {
	return func(r *requestInfo) {
		r.request.Header.SetCookie(k, v)
	}
}

// OptCookies attaches a set of key, value cookies to the request.
func OptCookies(c map[string]string) func(*requestInfo) {
	return func(r *requestInfo) {
		for k, v := range c {
			r.request.Header.SetCookie(k, v)
		}
	}
}

// OptCookies sets the request timeout (default if unset: 1 minute).
func OptTimeout(d time.Duration) func(*requestInfo) {
	return func(r *requestInfo) {
		r.timeout = d
	}
}
