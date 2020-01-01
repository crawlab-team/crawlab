package entity

type ConfigSpiderData struct {
	Version    string            `yaml:"version" json:"version"`
	Engine     string            `yaml:"engine" json:"engine"`
	StartUrl   string            `yaml:"start_url" json:"start_url"`
	StartStage string            `yaml:"start_stage" json:"start_stage"`
	Stages     []Stage           `yaml:"stages" json:"stages"`
	Settings   map[string]string `yaml:"settings" json:"settings"`
}

type Stage struct {
	Name      string  `yaml:"name" json:"name"`
	IsList    bool    `yaml:"is_list" json:"is_list"`
	ListCss   string  `yaml:"list_css" json:"list_css"`
	ListXpath string  `yaml:"list_xpath" json:"list_xpath"`
	PageCss   string  `yaml:"page_css" json:"page_css"`
	PageXpath string  `yaml:"page_xpath" json:"page_xpath"`
	PageAttr  string  `yaml:"page_attr" json:"page_attr"`
	Fields    []Field `yaml:"fields" json:"fields"`
}

type Field struct {
	Name      string `yaml:"name" json:"name"`
	Css       string `yaml:"css" json:"css"`
	Xpath     string `yaml:"xpath" json:"xpath"`
	Attr      string `yaml:"attr" json:"attr"`
	NextStage string `yaml:"next_stage" json:"next_stage"`
	Remark    string `yaml:"remark" json:"remark"`
}
