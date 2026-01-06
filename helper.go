package httpmatter

// matterer is a generic interface that can be used to create a matter
type matterer interface {
	WithOptions(opts ...Option) error
	Read() error
	Parse() error
}

// Request returns a http request with frontmatter for a given namespace and name
func Request(namespace, name string, opts ...Option) (*RequestMatter, error) {
	matter := NewRequestMatter(namespace, name)
	return matter, makeMatter(matter, opts...)
}

// Response returns a http response with frontmatter for a given namespace and name
func Response(namespace, name string, opts ...Option) (*ResponseMatter, error) {
	matter := NewResponseMatter(namespace, name)
	return matter, makeMatter(matter, opts...)
}

// makeMatter makes a matter with the given options
func makeMatter(matter matterer, opts ...Option) error {
	if err := matter.WithOptions(opts...); err != nil {
		return err
	}
	if err := matter.Read(); err != nil {
		return err
	}
	if err := matter.Parse(); err != nil {
		return err
	}
	return nil
}
