package server

import (
	"net/http"
	"strings"
)

type (
	// HandlerOptions hold information of a HTTP handler options.
	HandlerOptions struct {
		// mandatory
		p string
		h http.Handler

		// options
		m            []string
		q            []string
		hdr          []string
		prefix       bool
		interceptors []HTTPInterceptor
	}

	handlerOptionsSlice []HandlerOptions

	// HTTPInterceptor is an interceptor/middleware func.
	HTTPInterceptor = func(http.Handler) http.Handler
)

// NewHandlerOptions return new empty HTTP options.
func NewHandlerOptions() *HandlerOptions {
	return &HandlerOptions{}
}

// Prefix mark that the HTTP handler is a prefix handler.
func (r *HandlerOptions) Prefix() *HandlerOptions {
	r.prefix = true
	return r
}

// Methods adds a matcher for HTTP methods. It accepts a sequence of one or more methods to be matched.
func (r *HandlerOptions) Methods(methods ...string) *HandlerOptions {
	r.m = methods
	return r
}

// Queries adds a matcher for URL query values. It accepts a sequence of key/value pairs. Values may define variables.
func (r *HandlerOptions) Queries(queries ...string) *HandlerOptions {
	r.q = queries
	return r
}

// Headers adds a matcher for request header values. It accepts a sequence of key/value pairs to be matched.
func (r *HandlerOptions) Headers(headers ...string) *HandlerOptions {
	r.hdr = headers
	return r
}

// Interceptors adds interceptors into the handler.
func (r *HandlerOptions) Interceptors(interceptors ...HTTPInterceptor) *HandlerOptions {
	r.interceptors = interceptors
	return r
}

func (p handlerOptionsSlice) Len() int { return len(p) }

func (p handlerOptionsSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p handlerOptionsSlice) Less(i, j int) bool {
	if p[i].prefix && !p[j].prefix {
		return true
	}
	v := len(strings.Split(p[i].p, "/")) - len(strings.Split(p[j].p, "/"))
	if v != 0 {
		return v < 0
	}
	v = strings.Compare(p[i].p, p[j].p)
	if v != 0 {
		return v < 0
	}
	v = len(p[i].m) - len(p[j].m)
	if v != 0 {
		return v < 0
	}
	v = len(p[i].q) - len(p[j].q)
	if v != 0 {
		return v < 0
	}
	return len(p[i].hdr)-len(p[j].hdr) < 0
}
