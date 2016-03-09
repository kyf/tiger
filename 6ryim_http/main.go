package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fvbock/endless"
	"github.com/go-martini/martini"
)

const (
	LOG_PREFIX string = "[6ryim_http]"
	LOG_PATH   string = "/var/log/6ryim_http/6ryim_http.log"
)

var (
	configFile *string
	port       *string

	uri string
)

func init() {
	configFile = flag.String("config", "./conf.ini", "http service config file")
}

func response(w http.ResponseWriter, result interface{}) {
	data := make(map[string]interface{})
	switch r := result.(type) {
	case string:
		data["status"] = "error"
		data["msg"] = r
	case map[string]string:
		data["status"] = "ok"
		data["msg"] = "success"
		for k, v := range r {
			data[k] = v
		}
	default:
		data["status"] = "ok"
		data["msg"] = "success"
		data["data"] = result
	}
	re, _ := json.Marshal(data)
	w.Write(re)
}

var (
	handlers map[string]martini.Handler = make(map[string]martini.Handler)
	C        Config
)

func serveHTTP(context martini.Context, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	context.Map(r.Form)
	w.Header().Add("Access-Control-Allow-Origin", "*")
}

func main() {
	flag.Parse()

	fp, err := os.OpenFile(LOG_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fp)
	}
	defer fp.Close()

	logger := log.New(fp, LOG_PREFIX, log.LstdFlags)
	m := martini.Classic()
	m.Map(logger)
	m.Use(serveHTTP)

	C, err = initConfig(*configFile)
	if err != nil {
		logger.Printf("config error:%v", err)
		os.Exit(1)
	}

	for p, h := range handlers {
		m.Get(p, h)
		m.Post(p, h)
	}
	uri = fmt.Sprintf(":%s", C.port)
	var exit chan error
	go func() {
		exit <- endless.ListenAndServe(uri, m)
	}()

	err = <-exit
	logger.Printf("service exit, err is %v", err)
	os.Exit(1)
}
