package errhandler

type notFoundErr struct {
	inner error
}

func NewNotFoundErr(err error) error {
	return notFoundErr{
		inner: err,
	}
}

func (e notFoundErr) Error() string {
	return e.inner.Error()
}
