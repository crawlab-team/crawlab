package errors

func NewGitError(msg string) (err error) {
	return NewError(ErrorPrefixGit, msg)
}

var (
	ErrorGitInvalidAuthType = NewGitError("invalid auth type")
	ErrorGitNoMainBranch    = NewGitError("no main branch")
)
