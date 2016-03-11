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
	LOG_PREFIX string = "[6ryim_admin]"
)

var (
	LogPath string = "/var/log/6ryim_admin/6ryim_admin.log"
	Port    int
)

func init() {
	flag.IntVar(&Port, "port", 6060, "listen port")
}

func main() {
	m := martini.Classic()

	fp, err := os.OpenFile(LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("OpenFile failure, err is %v\n", err)
		os.Exit(1)
	}
	defer fp.Close()
	mylogger := log.New(fp, LOG_PREFIX, log.LstdFlags)
	m.Map(mylogger)

	m.Use(martini.Static("./static"))
	m.Use(martini.Static("./tpl"))

	var exit chan error
	go func() {
		exit <- endless.ListenAndServe(fmt.Sprintf(":%d", Port), m)
	}()

	e := <-exit
	mylogger.Printf("admin service exit err:%v", e)
	os.Exit(1)
}
