package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fvbock/endless"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	"bitbucket.org/kyf456/wx_service/user"
)

const (
	ORDER_REQUEST_URL = "http://admin.6renyou.com/weixin_plugin/getAjaxDetail?oid="

	LOG_PREFIX string = "[6ryim_admin]"

	TIME_LAYOUT string = "2006-01-02 15:04:05"

	SERVER_INVALID = "Server Invalid"
)

type m6ryResponse struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}

var (
	LogPath string = "/var/log/6ryim_admin/6ryim_admin.log"
	Port    int

	um   *UserManager
	wxum *user.UserManager
)

func init() {
	flag.IntVar(&Port, "port", 6060, "listen port")
}

func auth(r *http.Request, ren render.Render, logger *log.Logger, sess sessions.Session) {
	r.ParseForm()
	var admin_user string = ""
	admin_user, _ = sess.Get("admin_user").(string)

	authlist := []string{"/login", "/checklogin", "/request/receive", "/request/message/show"}
	extlist := []string{"css", "js", "jpg", "gif", "png"}
	ext := path.Ext(r.RequestURI)
	if !StringSliceContains(r.URL.Path, authlist) && !StringSliceContains(ext, extlist) && strings.EqualFold("", admin_user) {
		ren.Redirect("/login")
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
	m.Use(martini.Static(UPLOAD_PATH))

	m.Use(auth)

	um = NewUserManager()
	wxum = user.New(mylogger)

	m.Get("/message", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "index", data)
	})

	m.Get("/", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		ren.Redirect("/message")
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

		admin_user, _ := sess.Get("admin_user").(string)
		top, _ := fetchTop(admin_user)
		left, _ := fetchLeft()
		data := struct {
			Order string
			Left  template.HTML
			Top   template.HTML
		}{
			string(body),
			template.HTML(string(left)),
			template.HTML(string(top)),
		}
		ren.HTML(200, "message_detail", data)
	})

	m.Post("/wx/user/get", func(r *http.Request, ren render.Render, logger *log.Logger) {
		openids := r.Form.Get("openids")

		result := struct {
			m6ryResponse
			Data interface{} `json:"data"`
		}{}

		if len(openids) == 0 {
			result.Status = -1
			result.Info = "params is invalid!"
			ren.JSON(200, result)
			return
		}

		_openids := strings.Split(openids, ",")

		accessToken, err := getAccessToken()
		if err != nil {
			result.Status = -1
			result.Info = fmt.Sprintf("err:%v", err)
			logger.Printf("getAccessToken err:%v", err)
			ren.JSON(200, result)
			return
		}
		data := wxum.Get(_openids, accessToken)

		result.Status = 0
		result.Data = data
		ren.JSON(200, result)
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
		if admin_user, ok := sess.Get("admin_user").(string); ok && !strings.EqualFold("", admin_user) {
			defaultOL.release(admin_user)
		}
		sess.Delete("admin_user")
		ren.Redirect("/login")
	})

	m.Get("/login", func(logger *log.Logger, r *http.Request, sess sessions.Session, ren render.Render) {
		admin_user, ok := sess.Get("admin_user").(string)
		if ok && !strings.EqualFold("", admin_user) {
			ren.Redirect("/")
		} else {
			ren.HTML(200, "login", nil)
		}
	})

	m.Post("/checklogin", func(logger *log.Logger, r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user := r.Form.Get("user")
		admin_pwd := r.Form.Get("password")

		mgo := NewMongoClient()
		err := mgo.Connect()
		if err != nil {
			logger.Printf("mgo.Connect err:%v", err)
			ren.JSON(200, "failure")
			return
		}
		defer mgo.Close()

		adm := NewAdmin(admin_user, admin_pwd, "", mgo)
		status, err := adm.checkValid()
		if err != nil {
			logger.Printf("mgo.Connect err:%v", err)
			ren.JSON(200, "failure")
			return
		}

		if !status {
			ren.JSON(200, "failure")
			return
		}

		sess.Set("admin_user", admin_user)
		ren.JSON(200, "success")
	})

	m.Post("/admin/add", handleAdminAdd)
	m.Post("/admin/remove", handleAdminRemove)
	m.Post("/admin/list", handleAdminList)
	m.Get("/admin/list", handleAdminList)
	m.Post("/admin/edit", handleAdminEdit)

	m.Post("/request/receive", handleReceive)
	m.Post("/request/bind", handleBind)
	m.Post("/request/cc", handleRequestCC)
	m.Post("/request/send", handleSend)
	m.Post("/request/wait", handleListWait)
	m.Post("/request/fetch", handleFetchMsg)
	m.Post("/request/message/list", handleListMessage)
	m.Post("/request/message/show", handleListAllMessage)
	m.Post("/request/message/new/number", getNewMessageNum)
	m.Get("/call/center/message/detail", handleListDetail)
	m.Post("/upload", handleUpload)
	m.Post("/unbind", func(r *http.Request, w http.ResponseWriter, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		openid := r.Form.Get("openid")

		if strings.EqualFold("", openid) {
			responseJson(w, false, "openid is empty")
			return
		}
		defaultOL.unbind(admin_user, openid)
		responseJson(w, true, "")
	})
	m.Get("/request/online/list", handleOnlineList)
	m.Get("/request/cacheautoreply", handleCacheAutoReply)

	m.Post("/request/fastreply/add", handleFRadd)
	m.Post("/request/fastreply/update", handleFRupdate)
	m.Post("/request/fastreply/list", handleFRlist)
	m.Post("/request/fastreply/remove", handleFRremove)

	m.Post("/request/autoreply/timeitem/list", handleARlist)
	m.Post("/request/autoreply/timeitem/remove", handleARremove)
	m.Post("/request/autoreply/timeitem/update", handleARupdate)
	m.Post("/request/autoreply/timeitem/add", handleARadd)

	m.Get("/call/center/handled", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "handled", data)
	})

	m.Get("/call/center/fastreply", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "fastreply", data)
	})

	m.Get("/call/center/autoreply", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "autoreply", data)
	})

	m.Get("/call/center/account", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "account", data)
	})

	m.Get("/call/center/my/book", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "book", data)
	})

	m.Get("/call/center/message", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "ccmsglist", data)
	})

	m.Get("/call/center/wait", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "waitlist", data)
	})

	m.Get("/call/center/my", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}
		ren.HTML(200, "mycc", data)
	})

	m.Get("/chat", func(logger *log.Logger, r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		mgo := NewMongoClient()
		err := mgo.Connect()
		if err != nil {
			logger.Printf("mgo.Connect err:%v", err)
			ren.JSON(200, "failure")
			return
		}
		defer mgo.Close()

		adm, err := getAdminByName(admin_user, mgo)
		if err != nil {
			logger.Printf("getAdminByName err:%v", err)
			ren.JSON(200, "failure")
			return
		}

		data := struct {
			Opid template.HTML
			User template.HTML
		}{template.HTML(string(adm.Opid)), template.HTML(admin_user)}
		ren.HTML(200, "chat", data)
	})

	m.Get("/call/center/my/sendwx", func(r *http.Request, ren render.Render, sess sessions.Session) {
		admin_user, _ := sess.Get("admin_user").(string)
		left, _ := fetchLeft()
		top, _ := fetchTop(admin_user)
		data := struct {
			Left template.HTML
			Top  template.HTML
		}{template.HTML(string(left)), template.HTML(string(top))}

		ren.HTML(200, "sendwx", data)
	})

	handleCacheAutoReply(mylogger)

	var exit chan error
	go func() {
		exit <- endless.ListenAndServe(fmt.Sprintf(":%d", Port), m)
	}()

	e := <-exit
	mylogger.Printf("admin service exit err:%v", e)
	os.Exit(1)
}
