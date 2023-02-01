package main

import (
	"log"
	"time"
)

type Filter func(c *Context)

type FilterBuilder func(next Filter) Filter

// FilterReqLog 记录请求日志
func FilterReqLog(next Filter) Filter {
	return func(c *Context) {
		log.Println(c.r.URL)
		next(c)
	}
}

// FilterReqTime 记录请求时间
func FilterReqTime(next Filter) Filter {
	return func(c *Context) {
		s := time.Now().UnixNano()
		next(c)
		e := time.Now().UnixNano()
		log.Println(e - s)
	}
}
