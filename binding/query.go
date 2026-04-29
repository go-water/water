package binding

import "net/http"

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}

func (b queryBinding) Bind(req *http.Request, obj any) error {
	values := req.URL.Query()
	if err := b.mapTag(obj, values); err != nil {
		return err
	}

	return nil
}

func (b queryBinding) mapTag(ptr any, m map[string][]string) error {
	return mapFormByTag(ptr, m, b.Name())
}
