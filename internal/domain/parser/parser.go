package parser

type Package struct {
	Services []*Service
	Path     []string
	Key      string
}

type Service struct {
	Name      string
	Endpoints []*Endpoint
}

type Endpoint struct {
	Name   string
	Method string
	URL    string
	Path   []string
	Body   string
}
