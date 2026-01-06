package httpmatter

import (
	"io"
	"net/http"
	"net/http/httputil"
)

// ResponseMatter is a matter that can be used to store response content and error
type ResponseMatter struct {
	*Matter
	*http.Response
}

func NewResponseMatter(namespace, name string) *ResponseMatter {
	return &ResponseMatter{
		Matter: NewMatter(namespace, name),
	}
}

func (rm *ResponseMatter) Parse() error {
	content, err := rm.parse()
	if err != nil {
		return err
	}
	resp, err := ParseResponse(content)
	if err != nil {
		return err
	}
	rm.Response = resp
	return nil
}

func (rm *ResponseMatter) BodyString() (string, error) {
	body, err := rm.BodyBytes()
	if err != nil {
		return "", err
	}
	return string(body), nil
}
func (rm *ResponseMatter) BodyBytes() ([]byte, error) {
	body, err := io.ReadAll(rm.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (rm *ResponseMatter) Dump(resp *http.Response) error {
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	rm.content = string(b)
	rm.Response = resp
	return nil
}

func (rm *ResponseMatter) Save() error {
	return rm.Matter.Save()
}
