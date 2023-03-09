package web

import (
	"fmt"
)

type HandlerOnMap struct {
	handlers map[string]HandlerFunc
}

func NewHandlerOnMap() Handler {
	return &HandlerOnMap{handlers: make(map[string]HandlerFunc)}
}

func (h *HandlerOnMap) Route(method, pattern string, handler HandlerFunc, ms ...Middleware) {
	k := h.key(method, pattern)
	h.handlers[k] = handler
}

func (h *HandlerOnMap) key(method string, pattern string) string {
	return fmt.Sprintf("%s#%s", method, pattern)
}

func (h *HandlerOnMap) ServeHTTP(c *Context) {
	k := h.key(c.R.Method, c.R.URL.Path)
	if handler, ok := h.handlers[k]; ok {
		handler(c)
	} else {
		res := commonResp{
			Code: -1,
			Msg:  "路由不存在",
		}
		c.BadJson(res)
	}
}
