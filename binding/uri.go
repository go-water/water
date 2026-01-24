package binding

import (
	"net/http"
)

type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (uriBinding) Bind(req *http.Request, obj any) error {
	fields, ok := req.Context().Value(uriBinding{}).([]string)
	if !ok {
		return nil
	}

	m := make(map[string][]string, len(fields))
	for _, k := range fields {
		m[k] = []string{req.PathValue(k)}
	}

	if err := mapURI(obj, m); err != nil {
		return err
	}

	return validate(obj)
}
