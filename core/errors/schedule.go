package errors

func NewScheduleError(msg string) (err error) {
	return NewError(ErrorPrefixSchedule, msg)
}

//var ErrorSchedule = NewScheduleError("unregistered")
