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

func (c *Context) ReadJson(obj interface{}) error {
	all, err := io.ReadAll(c.r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(all, obj)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) WriteJson(code int, resp interface{}) error {
	c.w.WriteHeader(code)

	marshal, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = c.w.Write(marshal)
	return err
}

func (c *Context) OkJson(resp interface{}) error {
	err := c.WriteJson(http.StatusOK, resp)
	return err
}

func (c *Context) BadRequestJson(resp interface{}) error {
	err := c.WriteJson(http.StatusNotFound, resp)
	return err
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w, r}
}