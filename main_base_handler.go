package main

import (
	"fmt"
	"net/http"
)

type HandleBaseOnMap struct {
	// method + # + pattern
	router map[string]func(c *Context)
}

func NewHandleBaseOnMap() *HandleBaseOnMap {
	return &HandleBaseOnMap{router: make(map[string]func(c *Context))}
}

func (h *HandleBaseOnMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := h.key(r.Method, r.URL.Path)
	if handle, ok := h.router[key]; ok {
		handle(newContext(w, r))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

func (h *HandleBaseOnMap) key(method, pattern string) string {
	return fmt.Sprintf("%s#%s", method, pattern)
}
