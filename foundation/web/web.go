package web

import "net/http"

// BindQuery binds query params from URL to struct that implements QueryBinder
func BindQuery(r *http.Request, binder QueryBinder) error {
	return binder.Bind(QueryParams(r.URL.Query()))
}
