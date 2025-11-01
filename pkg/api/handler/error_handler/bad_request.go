package errhandler

type badRequestErr struct {
	inner error
}

func NewBadRequestErr(err error) error {
	return badRequestErr{
		inner: err,
	}
}

func (e badRequestErr) Error() string {
	return e.inner.Error()
}
