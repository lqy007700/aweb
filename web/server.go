package web

import (
	"context"
	"log"
	"net/http"
	"time"
)

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

func Sign(ctx *Context) {
	res := commonResp{
		Code: 0,
		Msg:  "success",
	}
	ctx.OkJson(res)
}

type Server interface {
	Routable
	Start()
	Shutdown(ctx context.Context) error
}

type webServer struct {
	name        string
	baseHandler Handler
	root        Filter
}

func (w *webServer) Route(method, pattern string, handler HandlerFunc, ms ...Middleware) {
	w.baseHandler.Route(method, pattern, handler, ms...)
}

func NewWebServer(name string, builders ...FilterBuilder) Server {
	// 将主函数定义为Filter放到链尾
	handler := NewHandlerOnTree()
	var root Filter = handler.ServeHTTP

	// 处理builders
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	return &webServer{
		name:        name,
		baseHandler: handler,
		root:        root,
	}
}

func (w *webServer) Start() {
	http.HandleFunc("/", func(wr http.ResponseWriter, r *http.Request) {
		ctx := NewContext(wr, r)
		// 一层层执行filter
		w.root(ctx)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func (w *webServer) Shutdown(ctx context.Context) error {
	log.Println(w.name, "退出中")
	time.Sleep(time.Second * 5)
	log.Println(w.name, "退出完毕")
	return nil
}
