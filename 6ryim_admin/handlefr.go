package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/martini-contrib/sessions"
	"gopkg.in/mgo.v2/bson"
)

func handleFRadd(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {
	content := r.Form.Get("content")

	if strings.EqualFold("", content) {
		responseJson(w, false, "content is empty!")
		return
	}

	mgo := NewMongoClient()
	err := mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	defer mgo.Close()
	admin_user, _ := sess.Get("admin_user").(string)
	fr := &FastReply{
		Id:         bson.NewObjectId(),
		Content:    content,
		Author:     admin_user,
		CreateTime: time.Now().Unix(),
	}

	err = fr.Add(mgo)
	if err != nil {
		logger.Printf("fastreply.Add err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "")
}

func handleFRupdate(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {
	content := r.Form.Get("content")
	id := r.Form.Get("id")

	if strings.EqualFold("", content) {
		responseJson(w, false, "content is empty!")
		return
	}

	if strings.EqualFold("", id) {
		responseJson(w, false, "id is empty!")
		return
	}

	mgo := NewMongoClient()
	err := mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	defer mgo.Close()
	fr := &FastReply{
		Id:      bson.ObjectIdHex(id),
		Content: content,
	}

	err = fr.Update(mgo)
	if err != nil {
		logger.Printf("fastreply.Update err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "")

}

func handleFRlist(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {
	mgo := NewMongoClient()
	err := mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	defer mgo.Close()
	data, err := FastReplyList(mgo)
	if err != nil {
		logger.Printf("FastReplyList err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "", data)
}

func handleFRremove(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {
	id := r.Form.Get("id")

	if strings.EqualFold("", id) {
		responseJson(w, false, "id is empty!")
		return
	}

	mgo := NewMongoClient()
	err := mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	defer mgo.Close()
	fr := &FastReply{
		Id: bson.ObjectIdHex(id),
	}

	err = fr.Remove(mgo)
	if err != nil {
		logger.Printf("fastreply.Remove err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "")

}
