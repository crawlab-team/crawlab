package errors

func NewSpiderError(msg string) (err error) {
	return NewError(ErrorPrefixSpider, msg)
}

var (
	ErrorSpiderMissingRequiredOption = NewSpiderError("missing required option")
	ErrorSpiderForbidden             = NewSpiderError("forbidden")
)
