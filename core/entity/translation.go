package entity

type Translation struct {
	Lang  string `json:"lang"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (t Translation) GetLang() (l string) {
	return t.Lang
}
