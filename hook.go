package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Hook func(c context.Context) error

func BuilderCloseServerHook(servers ...Server) Hook {
	return func(c context.Context) error {
		wg := sync.WaitGroup{}
		doneCh := make(chan struct{})

		wg.Add(len(servers))
		for _, s := range servers {
			go func(srv Server) {
				err := srv.Shutdown(c)
				if err != nil {
					log.Printf("server shutdown err: %v \n", err)
				}
				time.Sleep(time.Second)
				wg.Done()
			}(s)
		}

		// 正常退出监控
		go func() {
			wg.Wait()
			doneCh <- struct{}{}
		}()

		select {
		case <-c.Done(): // 超时
			log.Println("server shutdown time out")
			return errors.New("server shutdown time out")
		case <-doneCh: // 正常
			log.Println("server shutdown done")
			return nil
		}
	}
}

func BuilderNotifyGetwayHook() Hook {
	return func(c context.Context) error {
		fmt.Println("mock notify gateway")
		time.Sleep(time.Second * 2)
		return nil
	}
}
