package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kyf/postwx"
	"github.com/martini-contrib/sessions"
)

func handleReceive(r *http.Request, w http.ResponseWriter, logger *log.Logger) {
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

	switch MessageType(_msg_type) {
	case MSG_TYPE_IMAGE:
		fallthrough
	case MSG_TYPE_AUDIO:
		content, err = fetchWxMedia(content)
		if err != nil {
			logger.Printf("fetchWxMedia err:%v", err)
			responseJson(w, false, "media_id is invalid")
			return
		}
	}

	msg := Message{Fromtype: MSG_FROM_TYPE_USER, Openid: openid, Created: time.Now().Unix(), Content: content, MsgType: MessageType(_msg_type)}

	mgo := NewMongoClient()
	err = mgo.Connect()
	if err != nil {
		logger.Printf("mgo.Connect err:%v", err)
	} else {
		defer mgo.Close()
		err := storeMessage(msg, mgo)
		if err != nil {
			logger.Printf("storeMessage err:%v", err)
		}
	}

	if client := defaultOL.getClientByOpenid(openid); client != nil {
		client.appendMsg(msg)
	} else {
		defaultWL.Add(msg)
	}

	responseJson(w, true, "success")
}

func handleListMessage(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	openid := r.Form.Get("openid")

	if len(openid) == 0 {
		responseJson(w, false, "openid is empty")
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

	data, err := listMessage(openid, mgo)
	if err != nil {
		logger.Printf("listMessage err:%v", err)
		responseJson(w, false, SERVER_INVALID)
		return
	}

	responseJson(w, true, "", data)
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
	msg := ""
	if !status {
		msg = "接入失败"
	}
	responseJson(w, status, msg)
}

func handleFetchMsg(w http.ResponseWriter, r *http.Request, sess sessions.Session, logger *log.Logger) {
	openid := r.Form.Get("openid")
	if len(openid) == 0 {
		responseJson(w, false, "openid is empty")
		return
	}

	opid, _ := sess.Get("admin_user").(string)

	if len(opid) == 0 {
		responseJson(w, false, "opid is empty")
		return
	}

	defaultOL.poolLocker.Lock()
	defer defaultOL.poolLocker.Unlock()
	var result map[string]interface{}
	if clients, ok := defaultOL.olPool[opid]; ok {
		for index, client := range clients {
			if strings.EqualFold(client.openid, openid) {
				result = defaultOL.olPool[opid][index].fetchMsg()
				goto FOUND
			}
		}
	}

FOUND:

	responseJson(w, true, "", result)
}

func handleRequestCC(w http.ResponseWriter, r *http.Request, sess sessions.Session, logger *log.Logger) {
	opid, _ := sess.Get("admin_user").(string)

	if len(opid) == 0 {
		responseJson(w, false, "opid is empty")
		return
	}

	var data []map[string]interface{} = make([]map[string]interface{}, 0, 5)
	defaultOL.poolLocker.Lock()
	defer defaultOL.poolLocker.Unlock()
	if clients, ok := defaultOL.olPool[opid]; ok {
		for index, client := range clients {
			defaultOL.olPool[opid][index].refresh()

			msg := client.lastMsg.Content
			_ts := time.Unix(client.lastMsg.Created, 0)
			openid := client.openid
			times := _ts.Format(TIME_LAYOUT)
			ts := client.lastMsg.Created
			msgType := strconv.Itoa(int(client.lastMsg.MsgType))
			openid_name := openid
			isUpdate := false
			number := len(client.unRead)
			if len(client.unRead) > 0 {
				isUpdate = true
			}

			user, err := um.Get([]string{openid}, []string{"weixin"})
			if err != nil {
				logger.Printf("usermanager.Get err:%v", err)
			} else {
				if len(user) > 0 {
					openid_name = user[0].RealName
				}
			}

			data = append(data, map[string]interface{}{"number": number, "isUpdate": isUpdate, "msg": msg, "times": times, "ts": ts, "openid": openid, "msgType": msgType, "openid_name": openid_name})

		}
	}

	responseJson(w, true, "", data)

}

func handleSend(sess sessions.Session, w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	media_id, openid, message, msgType := r.Form.Get("media_id"), r.Form.Get("openid"), r.Form.Get("message"), r.Form.Get("msg_type")
	opid, _ := sess.Get("admin_user").(string)

	_msgType, err := strconv.Atoi(msgType)
	if err != nil {
		responseJson(w, false, fmt.Sprintf("msgtype [%v] is invalid!", msgType))
		return
	}

	var posterr error = nil

	switch MessageType(_msgType) {
	case MSG_TYPE_TEXT:
		_, posterr = postwx.PostText(openid, message)
	case MSG_TYPE_IMAGE:
		_, posterr = postwx.PostImage(openid, media_id)
	default:
		posterr = errors.New("Do not support wx message type!")
	}

	if posterr != nil {
		logger.Printf("postwx err:%v", posterr)
		responseJson(w, false, fmt.Sprintf("%v", posterr))
		return
	} else {
		msg := Message{Fromtype: MSG_FROM_TYPE_OP, Openid: openid, Created: time.Now().Unix(), Content: message, MsgType: MessageType(_msgType), Opid: opid}

		mgo := NewMongoClient()
		err = mgo.Connect()
		if err != nil {
			logger.Printf("mgo.Connect err:%v", err)
		} else {
			defer mgo.Close()
			err := storeMessage(msg, mgo)
			if err != nil {
				logger.Printf("storeMessage err:%v", err)
			}
		}
	}

	responseJson(w, true, "")

}

func handleListWait(w http.ResponseWriter) {
	data := make([]map[string]string, 0, len(defaultWL.waitPool))
	for openid, msgs := range defaultWL.waitPool {
		if len(msgs) == 0 {
			continue
		}
		msg := msgs[len(msgs)-1]
		_msg := msg.Content
		_ts := time.Unix(msg.Created, 0)
		ts := _ts.Format(TIME_LAYOUT)
		msgType := strconv.Itoa(int(msg.MsgType))

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

func uploadDir(prepath string) (string, error) {
	now := time.Now()
	path := fmt.Sprintf("%s/%v/%s/%v", prepath, now.Year(), now.Month(), now.Day())
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, nil
}

const (
	MAX_IMAGE_SIZE = 1024 * 1024 * 2
	UPLOAD_PATH    = "/mnt/uploads"
)

func handleUpload(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	var result map[string]string
	var cb string = "parent.callback"

	err := r.ParseMultipartForm(MAX_IMAGE_SIZE)
	if err != nil {
		cbResponseJson(w, false, fmt.Sprintf("err :%v", err), cb)
		return
	}

	if r.MultipartForm != nil && len(r.MultipartForm.File) > 0 {
		var msg string
		for _, files := range r.MultipartForm.File {
			if len(files) == 0 {
				msg = "no file upload"
				cbResponseJson(w, false, msg, cb)
				return
			}

			file := files[0]
			f, err := file.Open()
			if err != nil {
				logger.Printf("upload file err:%v", err)
				msg = "file error"
				cbResponseJson(w, false, msg, cb)
				return
			}

			d, err := ioutil.ReadAll(f)
			if err != nil {
				logger.Printf("upload file err:%v", err)
				msg = "read file error"
				cbResponseJson(w, false, msg, cb)
				return
			}

			data_len := int64(len(d))
			if data_len > MAX_IMAGE_SIZE {
				msg = "read file size more than maxImageSize"
				cbResponseJson(w, false, msg, cb)
				return
			}

			dir, err := uploadDir(UPLOAD_PATH)
			if err != nil {
				logger.Printf("uploadDir err:%v", err)
				msg = "upload dir error"
				cbResponseJson(w, false, msg, cb)
				return
			}

			if state := IsImage(file.Filename); !state {
				msg = "upload file extension invalid"
				cbResponseJson(w, false, msg, cb)
				return
			}

			fp := fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%v%s", time.Now().UnixNano(), filepath.Ext(file.Filename)))
			ioutil.WriteFile(fp, d, os.ModePerm)
			result = map[string]string{
				"filepath": strings.Replace(fp, UPLOAD_PATH, "", -1),
			}

			mediaType := "image"
			media_id, err := postwx.UploadMedia(fp, mediaType)
			if err == nil {
				result["media_id"] = media_id
			} else {
				logger.Printf("postwx.UploadMedia err:%v", err)
			}

			cbResponseJson(w, true, "", cb, result)
			break
		}
	} else {
		msg := "no file upload"
		cbResponseJson(w, false, msg, cb)
	}

}
