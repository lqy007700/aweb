package web

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type GracefulShutdown struct {
	reqCnt  int64 // 处理中的请求
	closing int32 // 关闭中的请求

	// 用 channel 来通知已经处理完了所有请求
	zeroReqCnt chan struct{}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		zeroReqCnt: make(chan struct{}),
	}
}

// ShutdownFilterBuilder 记录请求数和拒接请求
func (g *GracefulShutdown) ShutdownFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		loadInt32 := atomic.LoadInt32(&g.closing)

		// 拒绝所有请求
		if loadInt32 > 0 {
			log.Println("拒接请求")
			c.ShutdownJson()
			return
		}

		atomic.AddInt64(&g.reqCnt, 1)
		log.Printf("请求+1, %v", g)
		next(c)
		atomic.AddInt64(&g.reqCnt, -1)

		if g.closing > 0 && g.reqCnt == 0 {
			g.zeroReqCnt <- struct{}{}
		}
	}
}

func (g *GracefulShutdown) RejectNewRequestAndWaiting(c context.Context) error {
	// 拒绝请求
	log.Println("拒绝请求")
	atomic.AddInt32(&g.closing, 1)

	if atomic.LoadInt64(&g.reqCnt) == 0 {
		return nil
	}

	select {
	case <-c.Done():
		log.Println("超时")
		return errors.New("context 超时")
	case <-g.zeroReqCnt:
		log.Println("全部请求处理完成")
	}
	return nil
}
