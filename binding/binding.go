package binding

import "net/http"

type Binding interface {
	Name() string
	Bind(*http.Request, any) error
}

type StructValidator interface {
	ValidateStruct(any) error
	Engine() any
}

var Validator StructValidator = &defaultValidator{}

var (
	JSON = jsonBinding{}
	//XML   = xmlBinding{}
	//Form  = formBinding{}
	//Query = queryBinding{}
)

func validate(obj any) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
