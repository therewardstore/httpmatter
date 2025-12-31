package httpmatter

import (
	"io"
	"net/http"
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
	return body, nil
}
