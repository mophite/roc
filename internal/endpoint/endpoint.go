package endpoint

import (
	"roc/internal/namespace"
)

type Endpoint struct {
	Id      string
	Name    string
	Version string

	// schema.name.version.id
	// eg. goroc.api.hello.v.1.1.1.2d1bd2f9-6951-4235-83bd-d6f38b358552
	Absolute string
	Address  string

	// name.version
	// eg. api.hello.v.1.1.1
	Scope string
}

func (e *Endpoint) Splicing(schema string) *Endpoint {
	e.Scope = namespace.SplicingScope(e.Name, e.Version)
	e.Absolute = schema + "/" + e.Scope + "/" + e.Address
	return e
}
