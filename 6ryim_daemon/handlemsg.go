package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func handleMsg(msg *Message, logger *log.Logger) (bool, *Message) {
	//getMedia from weixin resource
	if strings.EqualFold(msg.ToType, TERMINAL_ADMIN) {
		if strings.EqualFold(msg.Source, MSG_SOURCE_WX) {
			switch msg.MsgType {
			case MSG_TYPE_IMAGE:
				fallthrough
			case MSG_TYPE_AUDIO:
				newpath, err := fetchWxMedia(msg.Message)
				if err != nil {
					logger.Printf("fetchWxMedia err:%v", err)
					break
				}
				msg.Message = newpath
			}
		}
	}

	if h.isOnline(msg.To) {
		return false, msg
	}

	var err error = nil
	err = storeOffline(*msg)
	if err != nil {
		logger.Printf("storeOffline err:%v", err)
	}

	if strings.EqualFold(msg.ToType, TERMINAL_ADMIN) {
		err = msg.sendAdminTpl()
	} else {

		switch msg.Source {
		case MSG_SOURCE_WX:
			err = msg.sendUserWX()
		case MSG_SOURCE_IOS:
			err = msg.sendUserIOS()
		case MSG_SOURCE_ANDROID:
			err = msg.sendUserAndroid()
		case MSG_SOURCE_360STREAM:
			err = msg.sendUser360Stream()
		default:
			return false, msg
		}
	}

	if err != nil {
		logger.Printf("handleMsg err:%v", err)
	}
	return true, msg
}

func fetchWxMedia(fpath string) (string, error) {
	data := make(url.Values)
	data.Set("fpath", fpath)
	res, err := http.PostForm(fmt.Sprintf("%swxmedia/fetch", HTTP_SERVICE_URL), data)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if newpath, ok := result["newpath"]; ok {
		return newpath, nil
	} else {
		return "", errors.New(result["msg"])
	}
}

func serveMsgReceive(w http.ResponseWriter, r *http.Request, logger *log.Logger, params url.Values) {
	m := params.Get("msg")

	if strings.EqualFold("", m) {
		response(w, false, "msg is empty")
		return
	}

	select {
	case h.message <- []byte(m):
		response(w, true, "success")
	default:
		response(w, false, "Server Invalid")
	}
}
