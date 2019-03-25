package config

type Config struct {
	Debug          bool                 `yaml:"debug"`
	DockerExecutor DockerExecutorConfig `yaml:"dockerExecutor"`
}

type DockerExecutorConfig struct {
	HostSocketPath string `yaml:"hostSocketPath"`
	Registry       string `yaml:"registry"`
	RegistryUser   string `yaml:"registryUser"`
	RegistryPass   string `yaml:"registryPass"`
}
