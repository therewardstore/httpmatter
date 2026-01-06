package httpmatter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertToGoTemplate(t *testing.T) {
	must := require.New(t)
	content := `{{who}}
	not a {{thing}},
	not a {{number}},
	not a {{date}}`
	out := convertToGoTemplate(content)
	must.Equal(`{{ index .Vars "who" }}
	not a {{ index .Vars "thing" }},
	not a {{ index .Vars "number" }},
	not a {{ index .Vars "date" }}`, out)
}

func TestExecuteTemplate(t *testing.T) {
	must := require.New(t)
	content := `{{ index .Vars "who" }}
	not a {{ index .Vars "thing" }},  \r\n
	not a {{ index .Vars "number" }},
	not a {{ index .Vars "date" }}`
	matter := &Matter{
		config: Config{},
		Vars: map[string]any{
			"who":    "John",
			"thing":  "table",
			"number": 123,
			"date":   "2021-01-01",
		},
	}
	out, err := executeTemplate(content, matter)
	must.Nil(err)
	must.Equal([]byte(`John
	not a table,  \r\n
	not a 123,
	not a 2021-01-01`), out)
}
