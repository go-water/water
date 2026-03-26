package binding

import (
	"encoding/json"
	"io"
	"net/http"
)

var EnableDecoderUseNumber = false
var EnableDecoderDisallowUnknownFields = false

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req.Body == nil || req.Body == http.NoBody || req.ContentLength == 0 {
		return nil
	}

	return decodeJSON(req.Body, obj)
}

func decodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return validate(obj)
}
