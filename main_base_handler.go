package main

import (
	"fmt"
	"net/http"
)

type Handle interface {
	ServeHTTP(c *Context)
	Routable
}

type HandleBaseOnMap struct {
	// method + # + pattern
	router map[string]func(c *Context)
}

func (h *HandleBaseOnMap) Route(method, pattern string, handle handleFunc) {
	key := h.key(method, pattern)
	h.router[key] = handle
}

func (h *HandleBaseOnMap) ServeHTTP(c *Context) {
	key := h.key(c.r.Method, c.r.URL.Path)
	if handle, ok := h.router[key]; ok {
		handle(newContext(c.w, c.r))
	} else {
		c.w.WriteHeader(http.StatusNotFound)
		_, _ = c.w.Write([]byte("not found"))
	}
}

func (h *HandleBaseOnMap) key(method, pattern string) string {
	return fmt.Sprintf("%s#%s", method, pattern)
}

func NewHandleBaseOnMap() Handle {
	return &HandleBaseOnMap{router: make(map[string]func(c *Context))}
}