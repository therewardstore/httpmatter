package httpmatter

import (
	"bufio"
	"strings"
	"testing"
)

// Matter is a generic matter that can be used to store content and error
type Matter struct {
	config    *Config
	front     string
	content   string
	Namespace string
	Name      string
	Vars      map[string]any
	tb        testing.TB
}

func NewMatter(namespace, name string) *Matter {
	return &Matter{
		config:    config,
		Namespace: namespace,
		Name:      name,
		Vars:      make(map[string]any),
	}
}

func (m *Matter) WithOptions(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(m); err != nil {
			return err
		}
	}
	return nil
}

func (m *Matter) Validate() error {
	_, err := openFile(m.filePath())
	if err != nil {
		return ErrReadingFile().WithData("file", m.filePath()).WithError(err)
	}
	return nil
}

// Read read the only request matter from the file
func (m *Matter) Read() error {
	// first read the .dot env file
	m.readDotEnv()
	m.ifTB(func(tb testing.TB) {
		tb.Log("Reading file", m.filePath(), "for", m.Namespace, m.Name)
	})
	front, content, err := readFile(m.filePath())
	if err != nil {
		return ErrReadingFile().WithData("file", m.filePath()).WithError(err)
	}
	// read content to matter.content
	m.front = front
	m.content = content
	return nil
}

// ReadOne read the request matter from the file and pick by name
func (m *Matter) ReadOne(name string) error {
	return ErrNotImplemented().WithData("method", "ReadOne").WithData("name", name)
}

func (m *Matter) parse() ([]byte, error) {
	out, err := executeTemplate(m.config.TemplateConverter(m.content), m)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// readDotEnv function will read the .dot env file and
// store the values in m.dotEnvs only if
// the file is found and read successfully
func (m *Matter) readDotEnv() {
	dotEnvPath := makeFilePath(
		m.config.BaseDir,
		m.Namespace,
		m.config.EnvFileName,
		m.config.EnvFileExtension)
	file, err := openFile(dotEnvPath)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if chunks := strings.SplitN(line, "=", 2); len(chunks) == 2 {
			m.Vars[chunks[0]] = chunks[1]
		}
	}
}

func (m *Matter) filePath() string {
	return makeFilePath(m.config.BaseDir, m.Namespace, m.Name, m.config.FileExtension)
}

func (m *Matter) ifTB(fn func(tb testing.TB)) {
	if m.tb == nil {
		return
	}
	fn(m.tb)
}
