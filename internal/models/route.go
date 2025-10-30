package models

type Route struct {
	method string
	path   string
}

func NewRoute(method, path string) *Route {
	return &Route{method: method, path: path}
}

func (r *Route) Method() string { return r.method }
func (r *Route) Path() string   { return r.path }
