package main

import (
	"context"
	"flag"
	"github.com/hanskorg/logkit"
	"os"
	"os/signal"
	"sms/conf"
	"sms/http"
	"sms/service"
	"syscall"
)

var (
	httpServer *http.Server
)

func main() {
	flag.Parse()
	conf.Init()
	httpServer = http.New(service.New())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			httpServer.Server.Shutdown(context.Background())
			logkit.Exit()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
