package app

type ModulePoll struct {
	Enabled   bool   `yaml:"enabled"`
	OwnServer bool   `yaml:"ownServer"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
}
