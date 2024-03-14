package render

import (
	"fmt"
	"github.com/go-water/water/internal/bytesconv"
	"net/http"
)

type Text struct {
	Format string
	Data   []any
}

var plainContentType = []string{"text/plain; charset=utf-8"}

func (r Text) Render(w http.ResponseWriter) error {
	return WriteString(w, r.Format, r.Data)
}

func (r Text) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, plainContentType)
}

func WriteString(w http.ResponseWriter, format string, data []any) (err error) {
	writeContentType(w, plainContentType)
	if len(data) > 0 {
		_, err = fmt.Fprintf(w, format, data...)
		return
	}
	_, err = w.Write(bytesconv.StringToBytes(format))
	return
}
