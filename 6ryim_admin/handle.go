package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/martini-contrib/sessions"
)

func handleReceive(r *http.Request, w http.ResponseWriter) {
	openid := r.Form.Get("openid")
	content := r.Form.Get("content")
	msgType := r.Form.Get("msgType")

	if len(openid) == 0 {
		responseJson(w, false, "openid is empty")
		return
	}

	if len(content) == 0 {
		responseJson(w, false, "content is empty")
		return
	}

	if len(msgType) == 0 {
		responseJson(w, false, "msgType is empty")
		return
	}

	_msg_type, err := strconv.Atoi(msgType)
	if err != nil {
		responseJson(w, false, "msgType is invalid")
		return
	}

	msg := Message{openid: openid, created: time.Now().Unix(), content: content, msgType: MessageType(_msg_type)}

	if client := defaultOL.getClientByOpenid(openid); client != nil {
		client.appendMsg(msg)
	} else {
		defaultWL.Add(msg)
	}

	responseJson(w, true, "success")
}

func handleBind(w http.ResponseWriter, r *http.Request, sess sessions.Session) {
	openid := r.Form.Get("openid")
	opid, _ := sess.Get("admin_user").(string)

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

func handleRequestCC(w http.ResponseWriter, r *http.Request, sess sessions.Session) {
	opid, _ := sess.Get("admin_user").(string)

	if len(opid) == 0 {
		responseJson(w, false, "opid is empty")
		return
	}

	clients := defaultOL.getClients(opid)
	var data []map[string]string = make([]map[string]string, 0, 5)
	for _, client := range clients {
		client.refresh()
		msg := client.lastMsg.content
		_ts := time.Unix(client.lastMsg.created, 0)
		openid := client.openid
		ts := _ts.Format(TIME_LAYOUT)
		msgType := strconv.Itoa(int(client.lastMsg.msgType))

		data = append(data, map[string]string{"msg": msg, "ts": ts, "openid": openid, "msgType": msgType})
	}

	responseJson(w, true, "", data)

}

func handleSend() {

}

func handleListWait(w http.ResponseWriter) {
	data := make([]map[string]string, 0, len(defaultWL.waitPool))
	for openid, msg := range defaultWL.waitPool {
		_msg := msg.content
		_ts := time.Unix(msg.created, 0)
		ts := _ts.Format(TIME_LAYOUT)
		msgType := strconv.Itoa(int(msg.msgType))

		data = append(data, map[string]string{"msg": _msg, "ts": ts, "openid": openid, "msgType": msgType})

	}
	responseJson(w, true, "", data)
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
