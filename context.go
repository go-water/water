package water

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-water/water/binding"
)

const ContextKey = "_go-water/context-key"

type H map[string]any

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	mu   sync.RWMutex
	Keys map[string]any
}

func (c *Context) Param(key string) string {
	return c.Request.PathValue(key)
}

func (c *Context) ShouldBindJSON(obj any) error {
	return c.ShouldBindWith(obj, binding.JSON)
}

func (c *Context) ShouldBindWith(obj any, b binding.Binding) error {
	return b.Bind(c.Request, obj)
}

func (c *Context) JSON(status int, data any) error {
	c.Writer.Header().Set("Content-Type", "application/json")

	if data == nil {
		c.Writer.WriteHeader(status)
		return nil
	}

	strJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.Writer.WriteHeader(status)
	_, err = c.Writer.Write(strJson)
	return err
}

func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
}

func (c *Context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.Keys[key]
	return
}

func BindJSON[T any](c *Context) (t *T, err error) {
	defer func() {
		if p := recover(); p != nil {
			switch pe := p.(type) {
			case error:
				err = pe
			default:
				err = fmt.Errorf("%s", pe)
			}
		}
	}()

	t = new(T)
	if err = c.ShouldBindJSON(t); err != nil {
		return nil, err
	}

	return t, nil
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// hasRequestContext returns whether c.Request has Context and fallback.
//func (c *Context) hasRequestContext() bool {
//	hasFallback := c.engine != nil && c.engine.ContextWithFallback
//	hasRequestContext := c.Request != nil && c.Request.Context() != nil
//	return hasFallback && hasRequestContext
//}

// Deadline returns that there is no deadline (ok==false) when c.Request has no Context.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	//if !c.hasRequestContext() {
	//	return
	//}
	return c.Request.Context().Deadline()
}

// Done returns nil (chan which will wait forever) when c.Request has no Context.
func (c *Context) Done() <-chan struct{} {
	//if !c.hasRequestContext() {
	//	return nil
	//}
	return c.Request.Context().Done()
}

// Err returns nil when c.Request has no Context.
func (c *Context) Err() error {
	//if !c.hasRequestContext() {
	//	return nil
	//}
	return c.Request.Context().Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key any) any {
	if key == 0 {
		return c.Request
	}
	if key == ContextKey {
		return c
	}
	if keyAsString, ok := key.(string); ok {
		if val, exists := c.Get(keyAsString); exists {
			return val
		}
	}
	//if !c.hasRequestContext() {
	//	return nil
	//}
	return c.Request.Context().Value(key)
}
