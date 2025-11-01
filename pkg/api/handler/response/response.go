package response

import (
	"encoding/json"
	"net/http"
)

type (
	Response[Out any] interface {
		OK(Out)
		NoContent()
	}

	response[Out any] struct {
		w http.ResponseWriter
	}
)

func New[Out any](w http.ResponseWriter) *response[Out] {
	return &response[Out]{
		w: w,
	}
}

func (r response[Out]) OK(data Out) {
	r.writeData(http.StatusOK, data)
}

func (r response[Out]) NoContent() {
	r.writeData(http.StatusNoContent, nil)
}

func (r response[Out]) BadRequest(err error) {
	r.writeErrorData(http.StatusBadRequest, err)
}

func (r response[Out]) NotFound(err error) {
	r.writeErrorData(http.StatusNotFound, err)
}

func (r response[Out]) InternalError(err error) {
	r.writeErrorData(http.StatusInternalServerError, err)
}

func (r response[Out]) writeErrorData(status int, err error) {
	r.writeData(status, errorDataResponse{
		Code:    status,
		Message: err.Error(),
	})
}

func (r response[Out]) writeData(status int, data any) {
	if data != nil {
		response, err := json.Marshal(data)
		if err != nil {
			r.writeInternalErrorText(err)
			return
		}

		r.w.Header().Set("Content-Type", "application/json")
		_, err = r.w.Write(response)
		if err != nil {
			r.writeInternalErrorText(err)
			return
		}
	}

	if status != http.StatusOK {
		r.w.WriteHeader(int(status))
	}
}

func (r response[Out]) writeInternalErrorText(err error) {
	r.w.WriteHeader(http.StatusInternalServerError)
	http.Error(r.w, err.Error(), http.StatusInternalServerError)
}
