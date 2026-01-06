package httpmatter

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
)

// RequestMatter is a matter that can be used to store request content and error
type RequestMatter struct {
	*Matter
	*http.Request
}

func NewRequestMatter(namespace, name string) *RequestMatter {
	return &RequestMatter{
		Matter: NewMatter(namespace, name),
	}
}

func (rm *RequestMatter) Parse() error {
	content, err := rm.parse()
	if err != nil {
		return err
	}
	req, err := ParseRequest(content)
	if err != nil {
		return err
	}
	rm.Request = req
	return nil
}

func (rm *RequestMatter) BodyString() (string, error) {
	body, err := rm.BodyBytes()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (rm *RequestMatter) BodyBytes() ([]byte, error) {
	body, err := io.ReadAll(rm.Body)
	if err != nil {
		return nil, err
	}
	// Preserve existing behavior (return the bytes) while also keeping the request
	// sendable after inspection by resetting the body reader.
	rm.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func (rm *RequestMatter) Dump(req *http.Request) error {
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	rm.content = string(b)
	rm.Request = req
	return nil
}

func (rm *RequestMatter) Save() error {
	return rm.Matter.Save()
}
