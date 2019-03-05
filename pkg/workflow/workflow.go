package workflow

type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Workflow struct {
	Name      string     `yaml:"name"`
	Variables []Variable `yaml:"variables"`
	Actions   []Action   `yaml:"actions"`
}
