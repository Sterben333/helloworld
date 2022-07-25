package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/library/conf/paladin.v2"
	"go-common/library/log"
	"go-common/library/net/trace"
	"helloworld/internal/di"

	_ "go.uber.org/automaxprocs"
)

func main() {
	flag.Parse()
	log.Init(nil) // debug flag: log.dir={path}
	defer log.Close()
	log.Info("helloworld start")
	trace.Init(nil)
	defer trace.Close()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info("helloworld exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
