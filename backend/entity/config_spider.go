package entity

type Field struct {
	Name  string `yaml:"name" json:"name"`
	Css   string `yaml:"css" json:"css"`
	Xpath string `yaml:"xpath" json:"xpath"`
	Attr  string `yaml:"attr" json:"attr"`
	Stage string `yaml:"stage" json:"stage"`
}

type Stage struct {
	List   bool    `yaml:"list" json:"list"`
	Css    string  `yaml:"css" json:"css"`
	Xpath  string  `yaml:"xpath" json:"xpath"`
	Fields []Field `yaml:"fields" json:"fields"`
}

type ConfigSpiderData struct {
	Version  string           `yaml:"version" json:"version"`
	StartUrl string           `yaml:"startUrl" json:"start_url"`
	Stages   map[string]Stage `yaml:"stages" json:"stages"`
}
