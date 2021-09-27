package web

// QueryBinder injects query params from URL to struct that implements this interface
type QueryBinder interface {
	Bind(params QueryParams) error
}
