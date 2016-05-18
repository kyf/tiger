package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kyf/postwx"
)

func fetchWxMedia(fpath string) (string, error) {
	data := make(url.Values)
	data.Set("fpath", fpath)
	dir, err := uploadDir(UPLOAD_PATH)
	if err != nil {
		return "", err
	}

	fp := fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%v", time.Now().UnixNano()))
	fullpath, err := postwx.GetMedia(fpath, fp)
	if err != nil {
		return "", err
	}
	newpath := strings.Replace(string(fullpath), UPLOAD_PATH, "", -1)
	return newpath, nil
}

func Pathinfo(file string) (map[string]string, error) {
	if len(file) == 0 {
		return nil, errors.New("file path is empty!")
	}
	var result map[string]string = make(map[string]string)
	info := strings.Split(file, "/")
	info_len := len(info)
	if info_len == 1 {
		result["basepath"] = ""
		result["filename"] = file
	} else {
		result["filename"] = info[info_len-1]
		result["basepath"] = strings.Replace(file, result["filename"], "", -1)
	}
	name := strings.Split(result["filename"], ".")
	result["extension"] = ""
	name_len := len(name)
	if name_len > 1 {
		result["extension"] = name[name_len-1]
	}

	return result, nil
}

func IsImage(file string) bool {
	pathinfo, err := Pathinfo(file)
	if err != nil {
		return false
	}
	exts := []string{"jpeg", "jpg", "png", "gif"}
	extension := strings.ToLower(pathinfo["extension"])
	for _, e := range exts {
		if strings.EqualFold(extension, e) {
			return true
		}
	}
	return false
}

func cbResponseJson(writer io.Writer, status bool, msg, cb string, data ...interface{}) {
	if len(msg) == 0 {
		msg = "success"
	}
	result := map[string]interface{}{
		"status": status,
		"msg":    msg,
	}

	if len(data) > 0 {
		result["data"] = data[0]
	}

	re, _ := json.Marshal(result)
	writer.Write([]byte(fmt.Sprintf("<script type='text/javascript'>%s(%s)</script>", cb, string(re))))
}

func responseJson(writer io.Writer, status bool, msg string, data ...interface{}) {
	if len(msg) == 0 {
		msg = "success"
	}
	result := map[string]interface{}{
		"status": status,
		"msg":    msg,
	}

	if len(data) > 0 {
		result["data"] = data[0]
	}

	re, _ := json.Marshal(result)
	writer.Write(re)
}

func fetchTop(admin_name string) ([]byte, error) {
	top, err := fetchFile("./tpl/top.html")
	if err != nil {
		return nil, err
	}
	_top := strings.Replace(string(top), "{admin_name}", admin_name, -1)
	return []byte(_top), nil
}

func fetchLeft() ([]byte, error) {
	return fetchFile("./tpl/left.html")
}

func fetchFile(path string) ([]byte, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	content, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func StringSliceContains(it string, its []string) bool {
	for _, item := range its {
		if strings.EqualFold(item, it) {
			return true
		}
	}

	return false
}

func getAccessToken() (string, error) {
	res, err := http.Get("http://m.6renyou.com/weixin_service/getAccessToken?account_type=1")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func postWeb(openid, message, msgType string) error {
	data := make(url.Values)
	data.Set("content", message)
	data.Set("openid", openid)
	data.Set("msgType", msgType)
	res, err := http.PostForm("http://localhost:6067/receive", data)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	response := new(struct {
		Status  bool   `json:"status"`
		Message string `json:"msg"`
	})

	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}

	if !response.Status {
		return errors.New(response.Message)
	}

	return nil
}

type CacheAutoReplyStruct struct {
	locker sync.RWMutex
	ar     []AutoReply
	far    []FirstAutoReply
}

var (
	CacheAutoReply = &CacheAutoReplyStruct{ar: make([]AutoReply, 0)}
)

func (car *CacheAutoReplyStruct) update(mgo *Mongo, logger *log.Logger) {
	car.locker.RLock()
	defer car.locker.RUnlock()

	ar, err := AutoReplyList(mgo)
	if err != nil {
		logger.Printf("AutoReplyList err:%v", err)
		return
	}

	far, err := FirstAutoReplyList(mgo)
	if err != nil {
		logger.Printf("FirstAutoReplyList err:%v", err)
		return
	}

	car.ar = ar
	car.far = far
}

func (car *CacheAutoReplyStruct) arlist() []AutoReply {
	return car.ar
}

func (car *CacheAutoReplyStruct) farlist() []FirstAutoReply {
	return car.far
}

func autoReply(openid string, source int, logger *log.Logger) {
	ar := CacheAutoReply.arlist()
	now := time.Now()
	year, month, day, location, st := now.Year(), now.Month(), now.Day(), now.Location(), now.Unix()

	for _, it := range ar {
		from := time.Date(year, month, day, it.FromHour, it.FromMinute, 0, 0, location)
		to := time.Date(year, month, day, it.ToHour, it.ToMinute, 0, 0, location)
		if st >= from.Unix() && st <= to.Unix() {
			var posterr error
			switch source {
			case MSG_SOURCE_WX:
				_, posterr = postwx.PostText(openid, it.Content)
			case MSG_SOURCE_PC:
				posterr = postWeb(openid, it.Content, fmt.Sprintf("%v", MSG_TYPE_TEXT))
			}

			if posterr != nil {
				logger.Printf("posterr is %v", posterr)
			} else {
				msg := Message{Fromtype: MSG_FROM_TYPE_OP, Openid: openid, Created: time.Now().Unix(), Content: it.Content, MsgType: MSG_TYPE_TEXT, Opid: SYSTEM}
				msg.Source = source
				mgo := NewMongoClient()
				err := mgo.Connect()
				if err != nil {
					logger.Printf("mgo.Connect err:%v", err)
				} else {
					defer mgo.Close()
					err = storeMessage(msg, mgo)
					if err != nil {
						logger.Printf("storeMessage err:%v", err)
					}
				}

			}
		}
	}
}

func welcome(openid string) error {
	far := CacheAutoReply.farlist()
	if far == nil || len(far) == 0 {
		return nil
	}
	posterr := postWeb(openid, far[0].Content, fmt.Sprintf("%v", MSG_TYPE_TEXT))
	if posterr != nil {
		return posterr
	} else {
		msg := Message{Fromtype: MSG_FROM_TYPE_OP, Openid: openid, Created: time.Now().Unix(), Content: far[0].Content, MsgType: MSG_TYPE_TEXT, Opid: SYSTEM}
		msg.Source = MSG_SOURCE_PC
		mgo := NewMongoClient()
		err := mgo.Connect()
		if err != nil {
			return err
		} else {
			defer mgo.Close()
			err = storeMessage(msg, mgo)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func filterHTML(content string) ([]byte, error) {
	reg, err := regexp.Compile("<[^>]+>")
	if err != nil {
		return nil, err
	}

	result := reg.ReplaceAll([]byte(content), []byte(""))
	return result, nil
}
