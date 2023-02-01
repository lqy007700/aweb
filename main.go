package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

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

func newWebServer(name string, builders ...FilterBuilder) Server {
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

func (w *webServer) Route(method string, pattern string, handler HandlerFunc) error {
	return w.baseHandler.Route(method, pattern, handler)
}

func (w *webServer) Shutdown(ctx context.Context) error {
	log.Println(w.name, "退出中")
	time.Sleep(time.Second * 5)
	log.Println(w.name, "退出完毕")
	return nil
}

func main() {
	shutdown := NewGracefulShutdown()
	server := newWebServer("app1", FilterReqTime, FilterReqLog, shutdown.ShutdownFilterBuilder)
	server.Route("GET", "/path/1/2", sign)
	go server.Start()

	WaitForShutdown(
		shutdown.RejectNewRequestAndWaiting, // 全部请求处理完后开始下线
		BuilderNotifyGetwayHook(),
		BuilderCloseServerHook(server),
	)
}
