package options

type LogOptions struct {
	Console  bool   `yaml:"Console"`
	Path     string `yaml:"Path"`
	LinkName string `yaml:"LinkName"`
}
