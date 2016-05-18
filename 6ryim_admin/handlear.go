package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/martini-contrib/sessions"
	"gopkg.in/mgo.v2/bson"
)

func handleARadd(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {
	content := r.Form.Get("content")
	fromHour := r.Form.Get("fromhour")
	fromMinute := r.Form.Get("fromminute")
	toHour := r.Form.Get("tohour")
	toMinute := r.Form.Get("tominute")

	if strings.EqualFold("", content) {
		responseJson(w, false, "content is empty!")
		return
	}

	if strings.EqualFold("", fromHour) {
		responseJson(w, false, "fromHour is empty!")
		return
	}

	if strings.EqualFold("", fromMinute) {
		responseJson(w, false, "fromminute is empty!")
		return
	}

	if strings.EqualFold("", toHour) {
		responseJson(w, false, "tohour is empty!")
		return
	}

	if strings.EqualFold("", toMinute) {
		responseJson(w, false, "tominute is empty!")
		return
	}

	_fromHour, err := strconv.Atoi(fromHour)
	if err != nil {
		responseJson(w, false, "fromhour is invalid!")
		return
	}

	_fromMinute, err := strconv.Atoi(fromMinute)
	if err != nil {
		responseJson(w, false, "fromminute is invalid!")
		return
	}

	_toHour, err := strconv.Atoi(toHour)
	if err != nil {
		responseJson(w, false, "tohour is invalid!")
		return
	}

	_toMinute, err := strconv.Atoi(toMinute)
	if err != nil {
		responseJson(w, false, "tominute is invalid!")
		return
	}

	mgo := NewMongoClient()
	err = mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	defer mgo.Close()
	ar := &AutoReply{
		Id:         bson.NewObjectId(),
		Content:    content,
		FromHour:   _fromHour,
		FromMinute: _fromMinute,
		ToHour:     _toHour,
		ToMinute:   _toMinute,
	}

	err = ar.Add(mgo)
	if err != nil {
		logger.Printf("autoreply.Add err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "")
}

func handleARupdate(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {
	content := r.Form.Get("content")
	id := r.Form.Get("id")
	fromHour := r.Form.Get("fromhour")
	fromMinute := r.Form.Get("fromminute")
	toHour := r.Form.Get("tohour")
	toMinute := r.Form.Get("tominute")

	if strings.EqualFold("", content) {
		responseJson(w, false, "content is empty!")
		return
	}

	if strings.EqualFold("", id) {
		responseJson(w, false, "id is empty!")
		return
	}

	if strings.EqualFold("", fromHour) {
		responseJson(w, false, "fromHour is empty!")
		return
	}

	if strings.EqualFold("", fromMinute) {
		responseJson(w, false, "fromminute is empty!")
		return
	}

	if strings.EqualFold("", toHour) {
		responseJson(w, false, "tohour is empty!")
		return
	}

	if strings.EqualFold("", toMinute) {
		responseJson(w, false, "tominute is empty!")
		return
	}

	_fromHour, err := strconv.Atoi(fromHour)
	if err != nil {
		responseJson(w, false, "fromhour is invalid!")
		return
	}

	_fromMinute, err := strconv.Atoi(fromMinute)
	if err != nil {
		responseJson(w, false, "fromminute is invalid!")
		return
	}

	_toHour, err := strconv.Atoi(toHour)
	if err != nil {
		responseJson(w, false, "tohour is invalid!")
		return
	}

	_toMinute, err := strconv.Atoi(toMinute)
	if err != nil {
		responseJson(w, false, "tominute is invalid!")
		return
	}

	mgo := NewMongoClient()
	err = mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	defer mgo.Close()
	ar := &AutoReply{
		Id:         bson.ObjectIdHex(id),
		Content:    content,
		FromHour:   _fromHour,
		FromMinute: _fromMinute,
		ToHour:     _toHour,
		ToMinute:   _toMinute,
	}

	err = ar.Update(mgo)
	if err != nil {
		logger.Printf("autoreply.Update err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "")

}

func handleARlist(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {
	mgo := NewMongoClient()
	err := mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	defer mgo.Close()
	data, err := AutoReplyList(mgo)
	if err != nil {
		logger.Printf("AutoReplyList err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "", data)
}

func handleARremove(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {
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
	ar := &AutoReply{
		Id: bson.ObjectIdHex(id),
	}

	err = ar.Remove(mgo)
	if err != nil {
		logger.Printf("autoreply.Remove err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "")

}

func handleARFirstLoad(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {

}

func handleARFirstSave(r *http.Request, logger *log.Logger, sess sessions.Session, w http.ResponseWriter) {

}
