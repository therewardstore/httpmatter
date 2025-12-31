package httpmatter

import "fmt"

var config *Config

type Config struct {
	BaseDir           string
	FileExtension     string
	EnvFileName       string
	EnvFileExtension  string
	TemplateConverter func(content string) string
}

func Init(conf *Config) error {
	if conf.BaseDir == "" {
		return fmt.Errorf("base dir is required")
	}
	// Default supported extensions is .http
	// Other valid values are .rest, .md etc
	if conf.FileExtension == "" {
		conf.FileExtension = ".http"
	}
	if conf.TemplateConverter == nil {
		conf.TemplateConverter = convertToGoTemplate
	}
	if conf.EnvFileName == "" {
		conf.EnvFileName = ""
	}
	if conf.EnvFileExtension == "" {
		conf.EnvFileExtension = ".env"
	}
	config = conf
	return nil
}
