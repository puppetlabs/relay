package config

type Config struct {
	Debug       bool   `yaml:"debug"`
	APIHostAddr string `yaml:"apiHostAddr"`
	CacheDir    string `yaml:"cacheDir"`
	TokenPath   string `yaml:"tokenPath"`
}
