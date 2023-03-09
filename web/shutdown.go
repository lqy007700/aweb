package web

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// ShutdownSignals receives shutdown signals to process
	ShutdownSignals = []os.Signal{
		os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM,
	}
)

// WaitForShutdown 优雅退出
func WaitForShutdown(hooks ...Hook) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdownSignals...)
	select {
	case sig := <-signals:
		log.Printf("get signal %s, application will shutdown \n", sig)
		// 超时2分钟 总时间
		time.AfterFunc(time.Minute*2, func() {
			log.Printf("get signal %s, application will shutdown \n", sig)
			os.Exit(1)
		})

		for _, hook := range hooks {
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*50)
			err := hook(ctx)
			if err != nil {
				log.Printf("退出错误 %v", err)
			}
			cancelFunc()
		}
		os.Exit(0)
	}
}
