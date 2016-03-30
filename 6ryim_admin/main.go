package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fvbock/endless"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

const (
	LOG_PREFIX string = "[6ryim_admin]"

	ADMIN_USER string = "6renyou"
	ADMIN_PWD  string = "6renyou.com"
)

var (
	LogPath string = "/var/log/6ryim_admin/6ryim_admin.log"
	Port    int
)

func init() {
	flag.IntVar(&Port, "port", 6060, "listen port")
}

func auth(r *http.Request, ren render.Render, logger *log.Logger, sess sessions.Session) {
	r.ParseForm()
	admin_user, ok := sess.Get("admin_user").(string)

	authlist := []string{"/", "/main", "/my/test", "/message/detail"}
	for _, it := range authlist {
		if strings.EqualFold(it, r.URL.Path) {
			if !ok || !strings.EqualFold(ADMIN_USER, admin_user) {
				ren.Redirect("/login")
			}
		}
	}

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

	store := sessions.NewCookieStore([]byte(LOG_PREFIX))

	m.Use(sessions.Sessions("admin_session", store))

	m.Use(render.Renderer(render.Options{Directory: "./tpl", Extensions: []string{".html"}}))
	m.Use(martini.Static("./static"))

	m.Use(auth)

	m.Get("/main", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		ren.HTML(200, "index", nil)
	})

	m.Get("/", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		ren.HTML(200, "index", nil)
	})

	m.Get("/my/test", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		ren.HTML(200, "demo", nil)
	})

	m.Get("/message/detail", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		ren.HTML(200, "message_detail", nil)
	})

	m.Get("/logout", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		sess.Delete("admin_user")
		ren.Redirect("/login")
	})

	m.Get("/login", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		admin_user, ok := sess.Get("admin_user").(string)
		if ok && strings.EqualFold(ADMIN_USER, admin_user) {
			ren.Redirect("/")
		} else {
			ren.HTML(200, "login", nil)
		}
	})

	m.Post("/checklogin", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user := r.Form.Get("user")
		admin_pwd := r.Form.Get("password")

		if strings.EqualFold(ADMIN_USER, admin_user) && strings.EqualFold(ADMIN_PWD, admin_pwd) {
			sess.Set("admin_user", admin_user)
			ren.JSON(200, "success")
		} else {
			ren.JSON(200, "failure")
		}
	})

	var exit chan error
	go func() {
		exit <- endless.ListenAndServe(fmt.Sprintf(":%d", Port), m)
	}()

	e := <-exit
	mylogger.Printf("admin service exit err:%v", e)
	os.Exit(1)
}
