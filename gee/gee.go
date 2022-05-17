package gee

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
}

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

//add router
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {

	engine.router.addRouter(method, pattern, handler)
}

//POST defines the method to add GET request
func (engine *Engine) GET(pattern string, handle HandlerFunc) {
	engine.addRouter("GET", pattern, handle)
}

//POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handle HandlerFunc) {
	engine.addRouter("POST", pattern, handle)
}

//implment ServeHTTP interface
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handler(c)
}

//Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
