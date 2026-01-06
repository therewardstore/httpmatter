# httpmatter
![Build Status](https://github.com/therewardstore/httpmatter/actions/workflows/ci.yml/badge.svg)

File-backed HTTP request/response fixtures with a thin wrapper on top of `github.com/jarcoal/httpmock`.

## Why / motive

`httpmatter` exists to make HTTP testing feel like working with real HTTP:

- Test **outgoing HTTP requests** with mocked responses.
- Store fixtures as **real HTTP messages** (request/response format).
- Go can parse this format easily via the standard library (`net/http`).
- IDE extensions like **REST Client** / **HttpYac** can run `.http` files directly from your editor.
- You can also load **incoming HTTP requests** from files and test handlers using `go test` (or your IDE).
- Improve **accessibility, readability, and testability** by keeping tests close to the actual HTTP.
- Cherry on top: **variables** for reusable fixtures.

## Install

```bash
go get github.com/therewardstore/httpmatter
```

## Fixture format

- Fixtures live under a `BaseDir/<namespace>/` directory.
- Files default to the `.http` extension.
- A file can have optional “front matter” (comments / metadata) before the HTTP message.
- Variables in the file may be referenced as `{{token}}` and will be substituted from `.Vars["token"]`.

Dotenv env files (optional):
- If configured, `EnvFileName` + `EnvFileExtension` (e.g. `.env.sample`) will be read from `BaseDir/<namespace>/`.
  - Example lookup: `BaseDir/<namespace>/<EnvFileName><EnvFileExtension>`
  - If `EnvFileName` is empty, it will look for: `BaseDir/<namespace>/<EnvFileExtension>` (e.g. `testdata/basic/.env.sample`)
- Format is `KEY=VALUE` (empty lines and `#` comments are ignored).
- Key/value pairs are merged into `.Vars`.

## Example fixture (`.http`)

This is a single HTTP request message with `{{vars}}` inside the HTTP message. The optional front matter is useful for IDE tools (REST Client / HttpYac).

```http
///
// @name create_order
@host=https://httpbin.org
@token=ExampleToken
///

POST {{host}}/post HTTP/1.1
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "ProductID": 42,
  "Quantity": 1
}
```

## Usage

### Load a response fixture

```go
package mypkg

import (
	"path/filepath"
	"testing"

	"github.com/therewardstore/httpmatter"
)

func init() {
	_ = httpmatter.Init(&httpmatter.Config{
		BaseDir: filepath.Join("testdata"),
		FileExtension: ".http",
	})
}

func TestSomething(t *testing.T) {
	resp, err := httpmatter.Response("basic", "response_with_header")
	if err != nil {
		t.Fatal(err)
	}

	body, err := resp.BodyString()
	if err != nil {
		t.Fatal(err)
	}
	_ = body
}
```

### Load and execute a request fixture

Load a `.http` fixture, substitute variables, and execute it using a standard `http.Client`.

```go
func TestRequestFixture(t *testing.T) {
	reqMatter, err := httpmatter.Request(
		"advanced",
		"create_order",
		httpmatter.WithVariables(map[string]any{
			"ProductID": 123,
			"token": "secret-token",
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(reqMatter.Request)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
```

### Capture and save a response fixture

Capture a real `*http.Response` and save it to a fixture file. This is useful for recording real API responses to use as future mocks.

```go
func TestRecordResponse(t *testing.T) {
	// Initialize the response matter. If the file doesn't exist yet,
	// Response() returns ErrReadingFile which we can ignore when recording.
	respMatter, err := httpmatter.Response("tmp", "recorded_api_response")
	if err != nil && !errors.Is(err, httpmatter.ErrReadingFile()) {
		t.Fatal(err)
	}

	// Make a real request using standard http.Client
	resp, err := http.Get("https://httpbin.org/json")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Capture the response content into the matter
	if err := respMatter.Dump(resp); err != nil {
		t.Fatal(err)
	}

	// Save to testdata/tmp/recorded_api_response.http
	if err := respMatter.Save(); err != nil {
		t.Fatal(err)
	}
}
```

### Mock outgoing HTTP calls (global)

This library uses `httpmock.Activate()` / `httpmock.DeactivateAndReset()`, which is **global within the current process**.

- Avoid `t.Parallel()` in tests that use `(*HTTP).Init()`.
```go

func init() {
	_ = httpmatter.Init(&httpmatter.Config{
		BaseDir: filepath.Join("testdata"),
		FileExtension: ".http",
	})
}

func TestVendorFlow(t *testing.T) {
	h := httpmatter.NewHTTP(t, "basic").
		Add("request_with_prompts_and_vars", "response_with_header").
		Respond(nil)

	h.Init()
	defer h.Destroy()

	// ... code under test that makes HTTP requests ...
}
```

## Limitations / notes

1. One file can contain only **one** HTTP request or **one** HTTP response.
2. Only `{{var}}` is supported for variable substitution **inside the HTTP message**.
   - For REST Client / HttpYac variable systems, use their own front matter/directives (like `@var=...`) for editor execution.
3. Since this package enables `httpmock` **globally** for outgoing requests, parallel tests in the same process are not supported.
   - Prefer running parallel **processes** (separate `go test` invocations) instead of `t.Parallel()`.

## License

MIT. See `LICENSE`.
