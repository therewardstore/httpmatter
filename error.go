package httpmatter

import (
	"fmt"
	"strings"
)

var ErrOpeningFile = newErrFn("failed to open file")
var ErrReadingFile = newErrFn("failed to read file")
var ErrParsingFile = newErrFn("failed to parse file")
var ErrParsingTemplate = newErrFn("failed to parse template")
var ErrExecutingTemplate = newErrFn("failed to execute template")
var ErrCreatingMatter = newErrFn("failed to create matter")
var ErrNotImplemented = newErrFn("not implemented")

type err struct {
	message string
	data    map[string]any
}

func newErrFn(message string) func() *err {
	return func() *err {
		return &err{
			message: message,
		}
	}
}

func (e *err) WithData(key string, value any) *err {
	if e.data == nil {
		e.data = make(map[string]any)
	}
	e.data[key] = value
	return e
}

func (e *err) Is(err error) bool {
	return strings.Contains(err.Error(), e.message+":")
}

func (e *err) WithError(err error) *err {
	e.WithData("error", err.Error())
	return e
}

func (e *err) Error() string {
	return fmt.Sprintf("%s: %+v", e.message, e.data)
}
