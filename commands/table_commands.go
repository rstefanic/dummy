package commands

type TableCommands struct {
	Name    string            `yaml:"name"`
	Count   int               `yaml:"count"`
	Columns map[string]string `yaml:"columns"`
}
