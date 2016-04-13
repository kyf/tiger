package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	ORDER_REQUEST_URL = "http://admin.6renyou.com/weixin_plugin/getAjaxDetail?oid="

	LOG_PREFIX string = "[6ryim_admin]"

	ADMIN_USER string = "6renyou"
	ADMIN_PWD  string = "6renyou.com"
)

type m6ryResponse struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}

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

	um := NewUserManager()

	m.Get("/main", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		ren.HTML(200, "index", nil)
	})

	m.Get("/", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		ren.HTML(200, "index", nil)
	})

	m.Get("/my/test", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		ren.HTML(200, "demo", nil)
	})

	m.Get("/message/detail", func(w http.ResponseWriter, logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		orderid := r.FormValue("orderid")
		if len(orderid) == 0 {
			fmt.Fprintf(w, "<div style='text-align:center;margin-top:150px;color:red;'>无效的订单</div>")
			return
		}

		res, err := http.Get(ORDER_REQUEST_URL + orderid)
		if err != nil {
			logger.Printf("request order detail err:%v", err)
			fmt.Fprintf(w, fmt.Sprintf("<div style='text-align:center;margin-top:150px;color:red;'>%v</div>", err))
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Printf("request err:%v", err)
			fmt.Fprintf(w, fmt.Sprintf("<div style='text-align:center;margin-top:150px;color:red;'>%v</div>", err))
			return
		}

		ren.HTML(200, "message_detail", string(body))
	})

	m.Post("/user/get", func(r *http.Request, ren render.Render, logger *log.Logger) {
		openids := r.Form.Get("openids")
		source := r.Form.Get("source")

		result := struct {
			m6ryResponse
			Data interface{} `json:"data"`
		}{}

		if len(openids) == 0 || len(source) == 0 {
			result.Status = -1
			result.Info = "params is invalid!"
			ren.JSON(200, result)
			return
		}

		_openids := strings.Split(openids, ",")
		_source := strings.Split(source, ",")

		if len(_openids) != len(_source) {
			result.Status = -1
			result.Info = "params is invalid!"
			ren.JSON(200, result)
			return
		}

		data, err := um.Get(_openids, _source)
		if err != nil {
			result.Status = -1
			result.Info = fmt.Sprintf("UserManager get err:%v", err)
			logger.Printf("UserManager get err:%v", err)
			ren.JSON(200, result)
			return
		}

		result.Status = 0
		result.Data = data
		ren.JSON(200, result)
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
