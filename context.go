package water

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-water/water/binding"
	"github.com/go-water/water/render"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	ContextKey = "_go-water/context-key"
)

var MaxMultipartMemory int64 = 32 << 20 // 32 MB

type H map[string]any

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	mu   sync.RWMutex
	Keys map[string]any

	sameSite http.SameSite
	wt       *Water

	queryCache url.Values
	formCache  url.Values
}

func (c *Context) reset() {
	c.sameSite = 0
	c.Keys = nil
	c.queryCache = nil
	c.formCache = nil
}

func (c *Context) Param(key string) string {
	return c.Request.PathValue(key)
}

func (c *Context) Query(key string) (value string) {
	value, _ = c.GetQuery(key)
	return
}

func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}

	return defaultValue
}

func (c *Context) QueryArray(key string) (values []string) {
	values, _ = c.GetQueryArray(key)
	return
}

func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.Request != nil {
			c.queryCache = c.Request.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

func (c *Context) PostForm(key string) (value string) {
	value, _ = c.GetPostForm(key)
	return
}

func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}

	return "", false
}

func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value, ok := c.GetPostForm(key); ok {
		return value
	}
	return defaultValue
}

func (c *Context) PostFormArray(key string) (values []string) {
	values, _ = c.GetPostFormArray(key)
	return
}

func (c *Context) GetPostFormArray(key string) (values []string, ok bool) {
	c.initFormCache()
	values, ok = c.formCache[key]
	return
}

func (c *Context) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		req := c.Request
		if err := req.ParseMultipartForm(c.wt.MaxMultipartMemory); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				fmt.Printf("error on parse multipart form array: %v", err)
			}
		}
		c.formCache = req.PostForm
	}
}

func (c *Context) ShouldBind(obj any) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

func (c *Context) ShouldBindJSON(obj any) error {
	return c.ShouldBindWith(obj, binding.JSON)
}

func (c *Context) ShouldBindQuery(obj any) error {
	return c.ShouldBindWith(obj, binding.Query)
}

func (c *Context) ShouldBindWith(obj any, b binding.Binding) error {
	err := b.Bind(c.Request, obj)
	switch err.(type) {
	case nil:
		return nil
	case validator.ValidationErrors:
		return err
	default:
		return Err(err.Error())
	}
}

func (c *Context) BindJSON(obj any) error {
	return c.MustBindWith(obj, binding.JSON)
}

func (c *Context) BindQuery(obj any) error {
	return c.MustBindWith(obj, binding.Query)
}

func (c *Context) BindHeader(obj any) error {
	return c.MustBindWith(obj, binding.Header)
}

func (c *Context) MustBindWith(obj any, b binding.Binding) error {
	if err := c.ShouldBindWith(obj, b); err != nil {
		return err
	}
	return nil
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

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *Context) GetStringMap(key string) (sm map[string]any) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

func (c *Context) Redirect(code int, location string) {
	c.Render(-1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.Request,
	})
}

func (c *Context) SetSameSite(sameSite http.SameSite) {
	c.sameSite = sameSite
}

func (c *Context) ContentType() string {
	content := c.GetHeader("Content-Type")
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}

	return content
}

func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: c.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) HTML(code int, name string, obj any) {
	instance := c.wt.HTMLRender.Instance(name, obj)
	c.Render(code, instance)
}

func (c *Context) Text(code int, format string, values ...any) {
	c.Render(code, render.Text{Format: format, Data: values})
}

func (c *Context) ClientIP() string {
	if c.wt.TrustedPlatform != "" {
		if addr := c.GetHeader(c.wt.TrustedPlatform); addr != "" {
			return addr
		}
	}

	remoteIP := net.ParseIP(c.RemoteIP())
	if remoteIP == nil {
		return ""
	}

	if c.wt.RemoteIPHeaders != nil {
		for _, headerName := range c.wt.RemoteIPHeaders {
			ip, valid := c.wt.validateHeader(c.GetHeader(headerName))
			if valid {
				return ip
			}
		}
	}
	return remoteIP.String()
}

func (c *Context) RemoteIP() string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

func (c *Context) Status(code int) {
	if code > 0 && code != http.StatusOK {
		c.Writer.WriteHeader(code)
	}
}

func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		return
	}

	if err := r.Render(c.Writer); err != nil {
		return
	}
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

	var obj any
	kind := reflect.TypeOf(t).Elem().Kind()
	if kind == reflect.Map {
		obj = reflect.MakeMap(reflect.TypeOf(t).Elem()).Interface()
	} else {
		obj = new(T)
	}

	if err = c.ShouldBind(obj); err != nil {
		return nil, err
	}

	if kind == reflect.Map {
		o := obj.(T)
		t = &o
	} else {
		return obj.(*T), nil
	}

	return t, nil
}

func (c *Context) hasRequestContext() bool {
	hasFallback := c.wt != nil && c.wt.ContextWithFallback
	hasRequestContext := c.Request != nil && c.Request.Context() != nil
	return hasFallback && hasRequestContext
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if !c.hasRequestContext() {
		return
	}
	return c.Request.Context().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	if !c.hasRequestContext() {
		return nil
	}
	return c.Request.Context().Done()
}

func (c *Context) Err() error {
	if !c.hasRequestContext() {
		return nil
	}
	return c.Request.Context().Err()
}

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
	if !c.hasRequestContext() {
		return nil
	}
	return c.Request.Context().Value(key)
}
