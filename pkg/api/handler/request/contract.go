package request

//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock_request

type (
	Params interface {
		Get(key string) (string, error)
	}

	Request[In any] interface {
		GetQuery() Params
		GetPath() Params
		GetBody() (In, error)
	}
)
