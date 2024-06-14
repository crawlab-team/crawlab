package errors

func NewStatsError(msg string) (err error) {
	return NewError(ErrorPrefixStats, msg)
}

var ErrorStatsInvalidType = NewStatsError("invalid type")
