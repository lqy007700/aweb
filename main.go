package main

import (
	"aweb/middleware/prometheus"
	"aweb/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// 一般来说，在实际中我们都会单独准备一个端口给这种监控
		http.ListenAndServe(":9090", nil)
	}()

	a := (&prometheus.MiddlewareBuilder{
		Subsystem: "web",
		Name:      "http_request",
		Help:      "这是测试例子",
		ConstLabels: map[string]string{
			"instance_id": "1234567",
		},
	}).Build()

	shutdown := web.NewGracefulShutdown()
	server := web.NewWebServer("app1", web.FilterReqTime, web.FilterReqLog, shutdown.ShutdownFilterBuilder)
	server.Route("GET", "/path", web.Sign, a)
	server.Start()
}
