package httpmatter

import (
	"maps"
	"testing"
)

type Option func(m *Matter) error

func WithVariables(vars map[string]any) Option {
	return func(m *Matter) error {
		if m.Vars == nil {
			m.Vars = make(map[string]any)
		}
		maps.Copy(m.Vars, vars)
		return nil
	}
}

func WithTB(tb testing.TB) Option {
	return func(m *Matter) error {
		m.tb = tb
		return nil
	}
}
