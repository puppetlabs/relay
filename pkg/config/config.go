package config

type Config struct {
	Debug       bool   `yaml:"debug"`
	APIHostAddr string `yaml:"apiHostAddr"`
	CachePath   string `yaml:"-"`
	TokenPath   string `yaml:"-"`
}
