package httpmatter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	Init(&Config{
		BaseDir:          filepath.Join(dir, "testdata"),
		EnvFileExtension: ".env.sample",
	})
}

func TestEnvFile(t *testing.T) {
	must := require.New(t)
	mock, err := Response("basic", "response_only_body")
	must.Nil(err)
	must.Len(mock.Vars, 2)
	must.Equal("https://httpbin.org", mock.Vars["host"])
	must.Equal("SuperSecretSerivces", mock.Vars["token"])
}

func TestMockOnlyBody(t *testing.T) {
	must := require.New(t)
	mock, err := Response("basic", "response_only_body")
	must.Nil(err)
	must.NotNil(mock.BodyBytes())
	must.Equal("response_only_body", mock.Name)
	must.Equal(200, mock.StatusCode)
	must.Equal("", mock.Header.Get("Content-Type"))
}

func TestMockWithHeader(t *testing.T) {
	must := require.New(t)
	mock, err := Response("basic", "response_with_header")
	must.Nil(err)
	must.NotNil(mock.BodyBytes())
	must.Equal("response_with_header", mock.Name)
	must.Equal(200, mock.StatusCode)
	must.Equal("application/json", mock.Header.Get("Content-Type"))
	must.Equal("Mon, 18 Aug 2025 13:44:32 GMT", mock.Header.Get("Date"))
	must.Equal("*", mock.Header.Get("Access-Control-Allow-Origin"))
}

func TestMockOnlyHeader(t *testing.T) {
	must := require.New(t)
	mock, err := Response("basic", "response_only_header")
	must.Nil(err)
	body, err := mock.BodyString()
	must.Nil(err)
	must.Equal("\r\n", body)
	must.Equal("response_only_header", mock.Name)
	must.Equal(200, mock.StatusCode)
	must.Equal("application/json", mock.Header.Get("Content-Type"))
	must.Equal("Mon, 18 Aug 2025 13:44:32 GMT", mock.Header.Get("Date"))
	must.Equal("*", mock.Header.Get("Access-Control-Allow-Origin"))
}

func TestMockWithPromptsAndVars(t *testing.T) {
	must := require.New(t)
	mock, err := Request(
		"advanced",
		"request_with_prompts_and_vars",
		WithVariables(map[string]any{
			"date": "2025-01-01T00:00:00Z",
			"user": "John Doe",
		}),
	)
	must.Nil(err)
	body, err := mock.BodyBytes()
	must.Nil(err)
	must.Len(body, 152)
	must.Equal("request_with_prompts_and_vars", mock.Name)
	must.Len(mock.Vars, 3+2) // 3 from the file, 2 from the options
	must.Equal("POST", mock.Method)
	must.Equal("https://httpbin.org/post", mock.URL.String())
	must.Equal("Bearer SuperSecretSerivces", mock.Header.Get("Authorization"))
	must.Equal("2025-01-01T00:00:00Z", mock.Header.Get("Date"))
	must.Equal("*", mock.Header.Get("Access-Control-Allow-Origin"))
}
