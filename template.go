package httpmatter

import (
	"bytes"
	"regexp"
	"text/template"
)

var indexVars = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)

// convertToGoTemplate http front matter is bit different from go template
// these are special cases this function will cover
// 1. {{<key>}} to {{index .Vars "<key>"}}
func convertToGoTemplate(content string) string {
	out := indexVars.ReplaceAllString(content, `{{ index .Vars "$1" }}`)
	return out
}

func executeTemplate(content string, matter *Matter) ([]byte, error) {
	tmpl, err := template.New(matter.filePath()).Parse(content)
	if err != nil {
		return nil, ErrParsingTemplate().WithError(err)
	}
	out := bytes.Buffer{}
	if err := tmpl.Execute(&out, matter); err != nil {
		return nil, ErrExecutingTemplate().WithError(err)
	}
	return out.Bytes(), nil
}
