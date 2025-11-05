package response

//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock_response

type (
	Success[Out any] interface {
		OK(Out)
		Raw([]byte) // 200
		NoContent()
	}

	Error interface {
		BadRequest(err error)
		NotFound(err error)
		TooManyRequests(err error)
		InternalError(err error)
	}
)
