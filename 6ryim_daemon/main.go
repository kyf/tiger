package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fvbock/endless"
	"github.com/go-martini/martini"
	"github.com/kyf/6ryim/util"
)

const (
	CERT_FILE string = "../certs/6ry.crt"
	KEY_FILE  string = "../certs/6ry.key"

	LOG_PATH   string = "/var/log/6ryim_daemon/6ryim_daemon.log"
	LOG_PREFIX string = "[6ryim_daemon]"

	HTTP_SERVICE_URL string = "http://127.0.0.1:8989/"
	PUSH_SERVICE_URL string = "http://im2.6renyou.com:3031/"

	PUSH_SERVICE_ACCESSID  string = "6renyou_20151222"
	PUSH_SERVICE_SECRETKEY string = "123456789"

	TERMINAL_ADMIN string = "1"
	TERMINAL_USER  string = "2"

	MSG_SOURCE_WX        string = "1"
	MSG_SOURCE_IOS       string = "2"
	MSG_SOURCE_ANDROID   string = "3"
	MSG_SOURCE_360STREAM string = "4"

	MSG_TYPE_TEXT  string = "2"
	MSG_TYPE_IMAGE string = "3"
	MSG_TYPE_AUDIO string = "4"

	MSG_SYSTEM_NAME string = "system"
	MSG_SYSTEM      string = "1"
	MSG_USER        string = "0"

	MSG_SYSTEM_TYPE_ORDER       string = "1"
	MSG_SYSTEM_TYPE_FETCH       string = "2"
	MSG_SYSTEM_TYPE_TRIP_SEND   string = "3"
	MSG_SYSTEM_TYPE_TRIP_SELECT string = "4"
	MSG_SYSTEM_TYPE_CANCEL      string = "5"
	MSG_SYSTEM_TYPE_ACTIVITY    string = "6"
	MSG_SYSTEM_TYPE_ERROR       string = "7"
)

var (
	Addr    string
	SslAddr string
)

func init() {
	flag.StringVar(&Addr, "port", "8060", "websocket daemon listen port")
	flag.StringVar(&SslAddr, "sslport", "4433", "websocket daemon listen port")
}

func main() {
	flag.Parse()
	m := martini.Classic()
	m.Use(auth)

	//fp, err := os.OpenFile(LOG_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	fp, err := util.NewWriter(LOG_PATH)
	if err != nil {
		fmt.Printf("OpenFile failure, err is %v", err)
		os.Exit(1)
	}
	defer fp.Close()
	mylog := log.New(fp, LOG_PREFIX, log.LstdFlags)
	go h.run(mylog)
	m.Map(mylog)
	m.Get("/:token", serveWS)
	m.Post("/message/receive", serveMsgReceive)
	m.Get("/online/list", listOnline)
	var exit chan error = make(chan error)
	go func() {
		exit <- endless.ListenAndServe(fmt.Sprintf(":%s", Addr), m)
	}()

	go func() {
		exit <- endless.ListenAndServeTLS(fmt.Sprintf(":%s", SslAddr), CERT_FILE, KEY_FILE, m)
	}()

	e := <-exit
	mylog.Printf("service exit:err is %v", e)
	os.Exit(1)
}
