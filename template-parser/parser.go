package parser

type Parser interface {
	Parse(template string) (err error)
	Render(args ...interface{}) (content string, err error)
}
