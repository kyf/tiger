package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"gopkg.in/mgo.v2/bson"
)

func handleListAllMessage(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	openid, _msg_type, keyword, _from_type := r.Form.Get("openid"), r.Form.Get("msg_type"), r.Form.Get("keyword"), r.Form.Get("fromtype")

	_size := r.Form.Get("size")
	_page := r.Form.Get("page")
	size := 20
	var err error
	if len(_size) > 0 {
		size, err = strconv.Atoi(_size)
		if err != nil {
			size = 20
		}
	}

	page := 1
	if len(_page) > 0 {
		page, err = strconv.Atoi(_page)
		if err != nil {
			page = 1
		}
	}

	msg_type := 0
	if len(_msg_type) > 0 {
		msg_type, _ = strconv.Atoi(_msg_type)
	}

	from_type := 0
	if len(_from_type) > 0 {
		from_type, _ = strconv.Atoi(_from_type)
	}

	where := make([]bson.M, 0, 2)
	if msg_type > 0 {
		where = append(where, bson.M{"msgtype": msg_type})
	}

	if len(keyword) > 0 {
		where = append(where, bson.M{"content": bson.M{"$regex": keyword}})
	}

	if len(openid) > 0 {
		where = append(where, bson.M{"openid": openid})
	}

	if from_type > 0 {
		where = append(where, bson.M{"fromtype": from_type})
	}

	var condition bson.M = nil
	if len(where) > 0 {
		if len(where) > 1 {
			condition = bson.M{"$and": where}
		} else {
			condition = where[0]
		}
	}

	skip := (page - 1) * size

	mgo := NewMongoClient()
	err = mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}
	defer mgo.Close()

	var sort_cond = "-_id"
	var result []Message
	err = mgo.Find(CC_MESSAGE_TABLE, condition).Sort(sort_cond).Skip(skip).Limit(size).All(&result)

	if err != nil {
		logger.Printf("search.list err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	var total int
	total, err = mgo.Find(CC_MESSAGE_TABLE, condition).Sort(sort_cond).Count()
	if err != nil {
		logger.Printf("search.total err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "", map[string]interface{}{"data": result, "total": total})
}

func getNewMessageNum(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	params := r.Form
	_lastid := params.Get("lastid")
	openid := params.Get("openid")
	_fromtype := params.Get("fromtype")

	cli := NewMongoClient()
	err := cli.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}
	defer cli.Close()

	where := make([]bson.M, 0)

	lastid := 0
	if len(_lastid) > 0 {
		lastid, err = strconv.Atoi(_lastid)
	}

	fromtype := 0
	if len(_fromtype) > 0 {
		fromtype, err = strconv.Atoi(_fromtype)
	}

	if lastid > 0 {
		where = append(where, bson.M{"created": bson.M{"$gt": lastid}})
	}

	if len(openid) > 0 {
		where = append(where, bson.M{"openid": openid})
	}

	if fromtype > 0 {
		where = append(where, bson.M{"fromtype": fromtype})
	}

	var condition bson.M = nil
	if len(where) > 0 {
		if len(where) > 1 {
			condition = bson.M{"$and": where}
		} else {
			condition = where[0]
		}
	}

	var total int
	total, err = cli.Find(CC_MESSAGE_TABLE, condition).Count()
	if err != nil {
		logger.Printf("cli.Find err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "", total)
}

func handleListDetail(sess sessions.Session, ren render.Render, w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	openid := r.FormValue("openid")
	if len(openid) == 0 {
		fmt.Fprintf(w, "<div style='text-align:center;margin-top:150px;color:red;'>无效的参数</div>")
		return
	}

	admin_user, _ := sess.Get("admin_user").(string)
	top, _ := fetchTop(admin_user)
	left, _ := fetchLeft()
	data := struct {
		Left template.HTML
		Top  template.HTML
	}{
		template.HTML(string(left)),
		template.HTML(string(top)),
	}
	ren.HTML(200, "cc_message_detail", data)

}
