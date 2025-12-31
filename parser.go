package httpmatter

import (
	"bufio"
	"bytes"
	"net/http"
	"strconv"
)

func ParseRequest(content []byte) (*http.Request, error) {
	normalizedContent := normalizeLineEndings(content)
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(normalizedContent)))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func ParseResponse(content []byte) (*http.Response, error) {
	// Normalize line endings to CRLF for HTTP parsing
	normalizedContent := normalizeLineEndings(content)
	response, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(normalizedContent)), nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// normalizeLineEndings normalizes line endings only in the header part (meta).
// The body (everything after the first blank line) is returned untouched.
func normalizeLineEndings(content []byte) []byte {
	// find header/body separator (try CRLFCRLF first, then LF LF, then CR CR)
	sepIdx := bytes.Index(content, []byte("\r\n\r\n"))
	sepLen := 4
	if sepIdx == -1 {
		if i := bytes.Index(content, []byte("\n\n")); i != -1 {
			sepIdx = i
			sepLen = 2
		} else if i := bytes.Index(content, []byte("\r\r")); i != -1 {
			sepIdx = i
			sepLen = 2
		}
	}

	var headerPart []byte
	var bodyPart []byte
	if sepIdx == -1 {
		// no explicit blank line -> everything is header (no body)
		headerPart = content
		bodyPart = nil
	} else {
		headerPart = content[:sepIdx]
		bodyPart = content[sepIdx+sepLen:]
	}

	// normalize header line endings:
	// 1) collapse CRLF -> LF, CR -> LF so we have a single separator
	tmp := bytes.ReplaceAll(headerPart, []byte("\r\n"), []byte("\n"))
	tmp = bytes.ReplaceAll(tmp, []byte("\r"), []byte("\n"))

	// 2) rebuild header with CRLF for each line
	lines := bytes.Split(tmp, []byte("\n"))
	var out bytes.Buffer
	for _, ln := range lines {
		// if line start with content-length,
		// skip it we will add it as last header
		if bytes.HasPrefix(ln, []byte("Content-Length:")) {
			continue
		}
		out.Write(ln)
		out.WriteString("\r\n")
	}

	// append body unchanged (no normalization in body)
	if len(bodyPart) == 0 {
		out.WriteString("\r\n")
	} else {
		out.WriteString("Content-Length: ")
		out.WriteString(strconv.Itoa(len(bodyPart)))
		out.WriteString("\r\n\r\n")
		out.Write(bodyPart)
	}

	return out.Bytes()
}
