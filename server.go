package main

import "time"

type signReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func sign(ctx *Context) {
	time.Sleep(10 * time.Second)
	res := commonResp{
		Code: 0,
		Msg:  "success",
	}
	ctx.OkJson(res)
}
