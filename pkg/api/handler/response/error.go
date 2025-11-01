package response

type errorDataResponse struct {
	// Code HTTP-код ошибки
	Code int `json:"code"`

	// Message Описание ошибки
	Message string `json:"message"`
}

type ErrorResponse interface {
	BadRequest(err error)
	NotFound(err error)
	InternalError(err error)
}
