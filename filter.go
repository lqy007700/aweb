package main

import (
	"fmt"
	"time"
)

type FilterBuilder func(next Filter) Filter

type Filter func(ctx *Context)

//var _ FilterBuilder = Builder

func Builder(next Filter) Filter {
	return func(ctx *Context) {
		s := time.Now().Nanosecond()
		next(ctx)
		e := time.Now().Nanosecond()
		fmt.Println(e - s)
	}
}

func Aprint(next Filter) Filter {
	return func(ctx *Context) {
		fmt.Println(1)
		next(ctx)
	}
}
