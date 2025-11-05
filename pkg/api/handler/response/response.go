package response

import (
	"encoding/json"
	"net/http"
)

type response[Out any] struct {
	w http.ResponseWriter
}

func New[Out any](w http.ResponseWriter) *response[Out] {
	return &response[Out]{
		w: w,
	}
}

func (r response[Out]) Raw(data []byte) {
	r.writeData(http.StatusOK, data)
}

func (r response[Out]) OK(data Out) {
	r.writeJSON(http.StatusOK, data)
}

func (r response[Out]) NoContent() {
	r.writeData(http.StatusNoContent, nil)
}

func (r response[Out]) BadRequest(err error) {
	r.writeError(http.StatusBadRequest, err)
}

func (r response[Out]) NotFound(err error) {
	r.writeError(http.StatusNotFound, err)
}

func (r response[Out]) TooManyRequests(err error) {
	r.writeError(http.StatusTooManyRequests, err)
}

func (r response[Out]) InternalError(err error) {
	r.writeError(http.StatusInternalServerError, err)
}

func (r response[Out]) writeError(status int, err error) {
	r.writeJSON(status, map[string]any{
		"code":    status,
		"message": err.Error(),
	})
}

func (r response[Out]) writeJSON(status int, data any) {
	response, err := json.Marshal(data)
	if err != nil {
		r.writeInternalErrorText(err)
		return
	}

	r.writeData(status, response)
}

func (r response[Out]) writeData(status int, data []byte) {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(status)

	if data == nil {
		return
	}

	_, err := r.w.Write(data)
	if err != nil {
		r.writeInternalErrorText(err)
	}
}

func (r response[Out]) writeInternalErrorText(err error) {
	r.w.WriteHeader(http.StatusInternalServerError)
	http.Error(r.w, err.Error(), http.StatusInternalServerError)
}
