package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fvbock/endless"
	"github.com/go-martini/martini"
	//"github.com/gorilla/websocket"
)

const (
	LOG_PATH   string = "/var/log/6ryim_daemon/6ryim_daemon.log"
	LOG_PREFIX string = "[6ryim_daemon]"
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

	go h.run()

	fp, err := os.OpenFile(LOG_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	mylog := log.New(fp, LOG_PREFIX, log.LstdFlags)
	m.Map(mylog)
	m.Get("/:token", serveWS)
	var exit chan error = make(chan error)
	go func() {
		exit <- endless.ListenAndServe(fmt.Sprintf(":%s", Addr), m)
	}()

	e := <-exit
	mylog.Printf("6ryim_daemon exit:err is %v", e)
	os.Exit(1)
}
