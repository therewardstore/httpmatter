package httpmatter

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLineEndingsContent(t *testing.T) {
	reqstr := `POST https://example.com/api/order HTTP/1.1
Content-Type: application/json
Authorization: Bearer {{token}}
Content-Length: 214

{
    "ProductID": 42,
    "Quantity": 1,
    "Denomination": 100,
    "RecipientFirstname": "John",
    "RecipientLastname": "Doe",
    "RecipientEmail": "john@doe.com",
    "CustomerOrderID": "COI123"
}`

	normalized := normalizeLineEndings([]byte(reqstr))

	should := assert.New(t)
	splits := bytes.Split(normalized, []byte("\r\n"))
	should.Equal(len(splits), 6)
	should.Equal(string(splits[0]), "POST https://example.com/api/order HTTP/1.1")
	should.Equal(string(splits[1]), "Content-Type: application/json")
	should.Equal(string(splits[2]), "Authorization: Bearer {{token}}")
	should.Equal(string(splits[3]), "Content-Length: 204")
	should.Equal(string(splits[4]), "")
	should.Equal(string(splits[5]), `{
    "ProductID": 42,
    "Quantity": 1,
    "Denomination": 100,
    "RecipientFirstname": "John",
    "RecipientLastname": "Doe",
    "RecipientEmail": "john@doe.com",
    "CustomerOrderID": "COI123"
}`)

}

func TestNormalizeLineEndingsWithoutContent(t *testing.T) {
	reqstr := `GET https://example.com/api/order HTTP/1.1
Authorization: Bearer {{token}}
Content-Length: 214`

	normalized := normalizeLineEndings([]byte(reqstr))

	should := assert.New(t)
	splits := bytes.Split(normalized, []byte("\r\n"))
	should.Len(splits, 4)
	should.Equal(string(splits[0]), "GET https://example.com/api/order HTTP/1.1")
	should.Equal(string(splits[1]), "Authorization: Bearer {{token}}")
	should.Equal(string(splits[2]), "")
	should.Equal(string(splits[3]), "")

}

func TestParseRequestPost(t *testing.T) {
	reqstr := `POST https://example.com/api/order HTTP/1.1
Content-Type: application/json
Authorization: Bearer {{token}}
Content-Length: 23(does not matter)

{
    "ProductID": 42,
    "Quantity": 1,
    "Denomination": 100,
    "RecipientFirstname": "John",
    "RecipientLastname": "Doe",
    "RecipientEmail": "john@doe.com",
    "CustomerOrderID": "COI123"
}`

	req, err := ParseRequest([]byte(reqstr))
	if err != nil {
		t.Fatal(err)
	}

	should := assert.New(t)

	should.Equal("https://example.com/api/order", req.URL.String())
	should.Equal("POST", req.Method)
	should.Equal("Bearer {{token}}", req.Header.Get("Authorization"))
	should.Equal("application/json", req.Header.Get("Content-Type"))
	should.NotNil(req.Body)
	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	should.Contains(string(body), `"ProductID": 42`)
	should.Contains(string(body), `"Quantity": 1`)
	should.Contains(string(body), `"Denomination": 100`)
	should.Contains(string(body), `"RecipientFirstname": "John"`)
	should.Contains(string(body), `"RecipientLastname": "Doe"`)
	should.Contains(string(body), `"RecipientEmail": "john@doe.com"`)
	should.Contains(string(body), `"CustomerOrderID": "COI123"`)
}

func TestParseRequestGet(t *testing.T) {
	reqstr := `GET https://example.com/api/order HTTP/1.1
Authorization: Bearer {{token}}
Content-Length: 23(does not matter)
`

	req, err := ParseRequest([]byte(reqstr))
	if err != nil {
		t.Fatal(err)
	}

	should := assert.New(t)

	should.Equal("https://example.com/api/order", req.URL.String())
	should.Equal("GET", req.Method)
	should.Equal("Bearer {{token}}", req.Header.Get("Authorization"))
	should.NotNil(req.Body)
	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	should.Equal("", string(body))
}

func TestParseResponse(t *testing.T) {
	respstr := `HTTP/1.1 200 OK
Server: nginx/1.24.0 (Ubuntu)
Content-Type: application/json
Transfer-Encoding: chunked
Connection: close
Cache-Control: no-cache, private
Date: Mon, 18 Aug 2025 13:44:32 GMT
Access-Control-Allow-Origin: *

{
  "code": 200,
  "status": "success",
  "response": {
    "productID": "88",
    "denomination": "1000",
    "available_quantity": 908,
    "stock_status": "available"
  }
}`
	resp, err := ParseResponse([]byte(respstr))
	if err != nil {
		t.Fatal(err)
	}

	should := assert.New(t)
	should.Equal(200, resp.StatusCode)
	should.Equal("application/json", resp.Header.Get("Content-Type"))
	should.Equal("no-cache, private", resp.Header.Get("Cache-Control"))
	should.Equal("Mon, 18 Aug 2025 13:44:32 GMT", resp.Header.Get("Date"))
	should.Equal("*", resp.Header.Get("Access-Control-Allow-Origin"))
}
