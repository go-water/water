package binding

import (
	"net/http"
)

type SetBindingKey struct{}

type setBinding struct{}

func (setBinding) Name() string {
	return "set"
}

func (setBinding) Bind(req *http.Request, obj any) error {
	m := make(map[string][]string)
	keys := make(map[string]any)
	ctx := req.Context()
	kv := ctx.Value(SetBindingKey{})
	switch pe := kv.(type) {
	case map[string]any:
		keys = pe
	default:
		keys = nil
	}

	for k, v := range keys {
		vv, _ := v.(string)
		m[k] = []string{vv}
	}

	if err := mapSet(obj, m); err != nil {
		return err
	}

	return nil
}
