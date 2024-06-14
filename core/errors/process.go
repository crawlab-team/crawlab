package errors

func NewProcessError(msg string) (err error) {
	return NewError(ErrorPrefixProcess, msg)
}

var (
	ErrorProcessReachedMaxErrors    = NewProcessError("reached max errors")
	ErrorProcessDaemonProcessExited = NewProcessError("daemon process exited")
)
