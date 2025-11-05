package test

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseRespBody[T any](resp *http.Response, v *T) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
