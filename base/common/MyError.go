package common

type BBSError struct {
	ErrorCode int32
	Message string
}

func (B *BBSError) Error() string {
	return B.Message
}

func NewError(Message string) *BBSError {
	return &BBSError{ErrorCode: FAIL, Message: Message}
}

func NewErrorWithCode(code int32, Message string) *BBSError {
	return &BBSError{ErrorCode: code, Message: Message}
}

func HandleError(err error) Result {
	bbsError,ok  := err.(*BBSError)
	if ok {
		return ErrorMeWithCode(bbsError.Error(), bbsError.ErrorCode)
	}
	return Error(err.Error())
}