package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fvbock/endless"
	"github.com/go-martini/martini"
)

const (
	LOG_PATH   string = "/var/log/6ryim_daemon/6ryim_daemon.log"
	LOG_PREFIX string = "[6ryim_daemon]"

	HTTP_SERVICE_URL string = "http://127.0.0.1:8989/"

	TERMINAL_ADMIN int = 1
	TERMINAL_USER  int = 2

	MSG_SOURCE_WX        int = 1
	MSG_SOURCE_IOS       int = 2
	MSG_SOURCE_ANDROID   int = 3
	MSG_SOURCE_360STREAM int = 4
)

var (
	Addr string
)

func init() {
	flag.StringVar(&Addr, "port", "8060", "websocket daemon listen port")
}

func main() {
	flag.Parse()
	m := martini.Classic()
	m.Use(auth)

	fp, err := os.OpenFile(LOG_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	mylog := log.New(fp, LOG_PREFIX, log.LstdFlags)
	go h.run(mylog)
	m.Map(mylog)
	m.Get("/:token", serveWS)
	var exit chan error = make(chan error)
	go func() {
		exit <- endless.ListenAndServe(fmt.Sprintf(":%s", Addr), m)
	}()

	e := <-exit
	mylog.Printf("service exit:err is %v", e)
	os.Exit(1)
}
