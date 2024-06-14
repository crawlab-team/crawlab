package parser

func Parse(template string, args ...interface{}) (content string, err error) {
	return ParseGeneral(template, args...)
}

func ParseGeneral(template string, args ...interface{}) (content string, err error) {
	p, _ := NewGeneralParser()
	if err := p.Parse(template); err != nil {
		return "", err
	}
	return p.Render(args...)
}
