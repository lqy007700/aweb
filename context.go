package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w: w, r: r}
}

func (c *Context) ReadJson(data interface{}) error {
	all, err := io.ReadAll(c.r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(all, data)
	return err
}

func (c *Context) WriteJson(code int, data interface{}) error {
	c.w.WriteHeader(code)
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.w.Write(marshal)
	return err
}

func (c *Context) OkJson(data interface{}) error {
	return c.WriteJson(http.StatusOK, data)
}

func (c *Context) BadJson(data interface{}) error {
	return c.WriteJson(http.StatusBadRequest, data)
}

func (c *Context) ShutdownJson() error {
	return c.WriteJson(http.StatusServiceUnavailable, nil)
}
