package ginx

type errorEntry struct {
	code  int
	msg   string
	cause error
}

func (e *errorEntry) Error() string {
	return e.msg
}

func (e errorEntry) GetCode() int {
	return e.code
}

func New(code int, msg string) error {
	return &errorEntry{code: code, msg: msg, cause: nil}
}

func Adapt(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*errorEntry); ok {
		return err
	}
	return &errorEntry{code: CodeUnknown, msg: err.Error(), cause: err}
}
