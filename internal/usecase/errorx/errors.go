package errorx

type ErrorX struct {
	Code string
	Err  error
}

func (e *ErrorX) Error() string {
	return e.Code
}

func (e *ErrorX) Unwrap() error {
	return e.Err
}

func New(code string, err error) *ErrorX {
	return &ErrorX{Code: code, Err: err}
}
