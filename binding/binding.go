package binding

import "net/http"

const (
	MIMEJSON              = "application/json"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

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
	JSON          = jsonBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormMultipart = formMultipartBinding{}
	Header        = headerBinding{}
	Uri           = uriBinding{}
)

func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return Form
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEMultipartPOSTForm:
		return FormMultipart
	default:
		return Form
	}
}

func validate(obj any) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
