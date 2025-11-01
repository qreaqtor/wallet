package request

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type (
	Request[In any] interface {
		GetQuery() query
		GetPath() path
		GetBody() (In, error)
	}

	request[In any] struct {
		req *http.Request
	}
)

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

func (r request[In]) GetPath() path {
	return path{
		data: mux.Vars(r.req),
	}
}

func (r request[In]) GetQuery() query {
	return query{
		data: r.req.URL.Query(),
	}
}

type path struct {
	data map[string]string
}

func (p path) Get(key string) (string, bool) {
	v, ok := p.data[key]
	return v, ok
}

type query struct {
	data url.Values
}

func (q query) Get(key string) string {
	return q.data.Get(key)
}
