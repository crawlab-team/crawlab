package errors

func NewDataSourceError(msg string) (err error) {
	return NewError(ErrorPrefixDataSource, msg)
}

var (
	ErrorDataSourceInvalidType           = NewDataSourceError("invalid type")
	ErrorDataSourceNotExists             = NewDataSourceError("not exists")
	ErrorDataSourceNotExistsInContext    = NewDataSourceError("not exists in context")
	ErrorDataSourceAlreadyExists         = NewDataSourceError("already exists")
	ErrorDataSourceMismatch              = NewDataSourceError("mismatch")
	ErrorDataSourceMissingRequiredFields = NewDataSourceError("missing required fields")
	ErrorDataSourceUnauthorized          = NewDataSourceError("unauthorized")
)
