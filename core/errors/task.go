package errors

func NewTaskError(msg string) (err error) {
	return NewError(ErrorPrefixTask, msg)
}

var (
	ErrorTaskNotExists     = NewTaskError("not exists")
	ErrorTaskAlreadyExists = NewTaskError("already exists")
	ErrorTaskNoNodeId      = NewTaskError("no node id")
	ErrorTaskNodeNotFound  = NewTaskError("node not found")
)
