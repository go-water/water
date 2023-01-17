package water

import (
	"fmt"
	"strings"
)

type ServerBase struct{}

func (s *ServerBase) Name(srv interface{}) string {
	fullName := fmt.Sprintf("%T", srv)
	index := strings.LastIndex(fullName, ".")
	name := fullName[index+1:]

	return name
}
