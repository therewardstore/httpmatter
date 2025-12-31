package httpmatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsContentLine(t *testing.T) {
	should := assert.New(t)
	should.True(isContentLine("GET /"))
	should.True(isContentLine("POST /"))
	should.True(isContentLine("PUT /"))
	should.True(isContentLine("DELETE /"))
	should.True(isContentLine("PATCH /"))
	should.True(isContentLine("HEAD /"))
	should.True(isContentLine("HTTP/1.1"))

	should.False(isContentLine("// This is a comment"))
	should.False(isContentLine("# This is a comment"))
	should.False(isContentLine(""))
}
