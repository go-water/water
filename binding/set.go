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
	bvMap := make(map[string]any)
	ctx := req.Context()
	o := ctx.Value(SetBindingKey{})
	switch ot := o.(type) {
	case map[string]any:
		bvMap = ot
	default:
		bvMap = nil
	}

	for k, v := range bvMap {
		val, _ := v.(string)
		m[k] = []string{val}
	}

	if err := mapSet(obj, m); err != nil {
		return err
	}

	return nil
}
