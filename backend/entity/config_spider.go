package entity

type Field struct {
	Name      string `yaml:"name" json:"name"`
	Css       string `yaml:"css" json:"css"`
	Xpath     string `yaml:"xpath" json:"xpath"`
	Attr      string `yaml:"attr" json:"attr"`
	NextStage string `yaml:"next_stage" json:"next_stage"`
}

type Stage struct {
	IsList  bool    `yaml:"is_list" json:"is_list"`
	ListCss string  `yaml:"list_css" json:"list_css"`
	PageCss string  `yaml:"page_css" json:"page_css"`
	Fields  []Field `yaml:"fields" json:"fields"`
}

type ConfigSpiderData struct {
	Version    string           `yaml:"version" json:"version"`
	Engine     string           `yaml:"engine" json:"engine"`
	StartUrl   string           `yaml:"start_url" json:"start_url"`
	StartStage string           `yaml:"start_stage" json:"start_stage"`
	Stages     map[string]Stage `yaml:"stages" json:"stages"`
}
