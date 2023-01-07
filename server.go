package main

import (
	"log"
)

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func signUp(c *Context) {
	req := &signUpReq{}
	resp := &commonResponse{
		Data: req,
	}

	err := c.ReadJson(req)
	if err != nil {
		resp.Msg = err.Error()
		c.BadRequestJson(resp)
		return
	}

	err = c.OkJson(resp)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
