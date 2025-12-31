package httpmatter

import (
	"net/http"
)

type requester func(req *http.Request) int
type responder func(req *http.Request, reqm *RequestMatter, respms []*ResponseMatter) *ResponseMatter

// DefaultResponder returns the first response without any assertions
var DefaultResponder responder = func(req *http.Request, reqm *RequestMatter, respms []*ResponseMatter) *ResponseMatter {
	return respms[0]
}

// RequestResponse returns a responder that asserts the request and returns the response at the index returned by the asserter
var RequestResponse = func(ca requester) responder {
	return func(req *http.Request, reqm *RequestMatter, respms []*ResponseMatter) *ResponseMatter {
		return respms[ca(req)]
	}
}
