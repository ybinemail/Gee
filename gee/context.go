package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request

	//request info
	Path   string
	Methon string

	//response info
	StatusCode int

	//params
	Params map[string]string

	//middleware
	handlers []HandlerFunc
	index    int
}

// construct new context
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Methon: r.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	handlerCount := len(c.handlers)
	for ; c.index < handlerCount; c.index++ {
		// do HandlerFunc func
		c.handlers[c.index](c)
	}
}

// set status
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// set header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// set string response
func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)

	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

//set json response
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
	//TODO json
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

//set data
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//set html
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

//get param by key
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// get postform by key
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//get Param
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
