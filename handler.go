package main

type handleFunc func(ctx *Context)

type Handle interface {
	ServeHTTP(c *Context)
	Routable
}