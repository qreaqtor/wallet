package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type request[In any] struct {
	req *http.Request
}

func New[In any](r *http.Request) Request[In] {
	return &request[In]{
		req: r,
	}
}

func (r request[In]) GetBody() (In, error) {
	var in In

	body, err := io.ReadAll(r.req.Body)
	if err != nil {
		return in, err
	}

	err = r.req.Body.Close()
	if err != nil {
		return in, err
	}

	err = json.Unmarshal(body, &in)
	if err != nil {
		return in, err
	}

	return in, nil
}

func (r request[In]) GetPath() Params {
	return path{
		data: mux.Vars(r.req),
	}
}

func (r request[In]) GetQuery() Params {
	return query{
		data: r.req.URL.Query(),
	}
}

type path struct {
	data map[string]string
}

func (p path) Get(key string) (string, error) {
	v, ok := p.data[key]

	if !ok {
		return "", fmt.Errorf("path param not found: %s", key)
	}

	return v, nil
}

type query struct {
	data url.Values
}

func (q query) Get(key string) (string, error) {
	v := q.data.Get(key)

	if v == "" {
		return "", fmt.Errorf("empty query param: %s", key)
	}

	return v, nil
}
