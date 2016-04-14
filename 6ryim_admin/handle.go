package main

import (
	"log"
	"net/http"
	//"time"
)

func handleReceive(r *http.Request, w http.ResponseWriter) {
	openid := r.Form.Get("openid")
	content := r.Form.Get("content")

	if len(openid) == 0 {
		responseJson(w, false, "openid is empty")
		return
	}

	if len(content) == 0 {
		responseJson(w, false, "content is empty")
		return
	}

	//msg := Message{openid: openid, created: time.Now().Unix(), content: content}

	/*
		if client := defaultOL.getClient(opid, openid); client != nil {
			client.appendMsg(msg)
			responseJson(w, true, "success")
			return
		} else {
			defaultWL.Add(msg)
			responseJson(w, true, "success")
			return
		}

	*/
}

func handleBind(w http.ResponseWriter, r *http.Request) {
	openid := r.Form.Get("openid")
	opid := r.Form.Get("opid")

	if len(openid) == 0 {
		responseJson(w, false, "openid is empty")
		return
	}

	if len(opid) == 0 {
		responseJson(w, false, "opid is empty")
		return
	}

	status := defaultWL.Fetch(opid, openid)
	responseJson(w, status, "")
}

func handleFetch() {

}

func handleSend() {

}

func handleAdminAdd(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	user, pwd, openid := r.FormValue("user"), r.FormValue("pwd"), r.FormValue("openid")

	if len(user) == 0 || len(pwd) == 0 || len(openid) == 0 {
		responseJson(w, false, "params is invalid")
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

	adm := NewAdmin(user, pwd, openid, mgo)
	ok, err := adm.checkUniq()
	if err != nil {
		logger.Printf("adm.checkUniq err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}
	if ok {
		err = adm.add()
		if err != nil {
			logger.Printf("adm.Add err:%v", err)
			responseJson(w, false, SERVER_INVALID)
			return
		} else {
			responseJson(w, true, "success")
			return
		}
	} else {
		responseJson(w, false, "该用户名已添加")
		return
	}

}

func handleAdminRemove(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	id := r.FormValue("id")

	if len(id) == 0 {
		responseJson(w, false, "params is invalid")
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

	adm := NewAdmin("", "", "", mgo)
	err = adm.remove(id)
	if err != nil {
		logger.Printf("adm.remove err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "")
}

func handleAdminList(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	mgo := NewMongoClient()
	err := mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	defer mgo.Close()

	adm := NewAdmin("", "", "", mgo)
	result, err := adm.list()
	if err != nil {
		logger.Printf("adm.list err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "success", result)
}

func handleAdminEdit(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	id, user, pwd, openid := r.FormValue("id"), r.FormValue("user"), r.FormValue("pwd"), r.FormValue("openid")

	if len(user) == 0 || len(pwd) == 0 || len(openid) == 0 || len(id) == 0 {
		responseJson(w, false, "params is invalid")
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

	adm := NewAdmin(user, pwd, openid, mgo)
	err = adm.edit(id)
	if err != nil {
		logger.Printf("adm.Edit err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "success")
}
