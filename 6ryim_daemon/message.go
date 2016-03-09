package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	im_type "github.com/kyf/6ryim/6ryim_http/im_type"
)

const (
	WX_TPL_ID    string            = "e5sFqp2BHA4OhbzOpzeqmi0ir6lT9sA3DanMOYOPhRI"
	WX_TPL_URL   string            = "http://m.6renyou.com/chat/index"
	WX_TPL_COLOR map[string]string = map[string]string{
		"top":    "#FF0000",
		"first":  "#000000",
		"data":   "#3eb166",
		"remark": "#939393",
	}

	WX_TPL string = `{
		"touser":"%s",
		"template_id":"%s",
		"url":"%s",
		"topcolor":"%s",
		"data":{
			"first":{
				"value":"",
				"color":"%s"
			},
			"keyword1":{
				"value":"%s",
				"color":"%s"
			},
			"keyword2": {
				"value":"%s",
				"color":"%s"
			},
			"keyword3": {
				"value":"%s",
				"color":"%s"
			},
			"remark":{
				"value":"",
				"color":"%s"
			}
		}
	}`

	WX_TPL_SOURCE_WEIXIN    string = "微信客户端"
	WX_TPL_SOURCE_IOS       string = "ios客户端"
	WX_TPL_SOURCE_ANDROID   string = "android客户端"
	WX_TPL_SOURCE_360STREAM string = "360stream客户端"
)

type Message im_type.Message

func newMsg(m []byte) (*Message, error) {
	var result Message
	err := json.Unmarshal(m, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *Message) sendAdminTpl() error {
	url := fmt.Sprintf("%s?id=%s", WX_TPL_URL, m.OrderId)
	date := getFormatNow("zh")

	var source string
	switch m.Source {
	case MSG_SOURCE_WX:
		source = WX_TPL_SOURCE_WEIXIN
	case MSG_SOURCE_IOS:
		source = WX_TPL_SOURCE_IOS
	case MSG_SOURCE_ANDROID:
		source = WX_TPL_SOURCE_ANDROID
	case MSG_SOURCE_360STREAM:
		source = WX_TPL_SOURCE_360STREAM
	}

	var msgbody string
	switch m.MsgType {
	case MSG_TYPE_TEXT:
		msgbody = m.Message
	default:
		msgbody = "您有一条新消息"
	}
	d := fmt.Sprintf(WX_TPL, m.To, WX_TPL_ID, url, WX_TPL_COLOR["top"], WX_TPL_COLOR["first"], source, WX_TPL_COLOR["data"], msgbody, WX_TPL_COLOR["data"], date, WX_TPL_COLOR["data"], WX_TPL_COLOR["remark"])
	_, err := postwx.PostTpl(d)
	return err
}

func (m *Message) sendUserWX() error {
	var posterr error = nil

	switch msg.MsgType {
	case MSG_TYPE_TEXT:
		_, posterr = postwx.PostText(m.To, m.Message)
	case MSG_TYPE_IMAGE:
		_, posterr = postwx.PostImage(m.To, m.Message)
	default:
		return errors.New("Do not support wx message type!")
	}

	if posterr != nil {
		var returnmsg Message = Message{
			From:       MSG_SYSTEM_NAME,
			To:         msg.From,
			Message:    fmt.Sprintf("%s:%v", msg.Message, posterr),
			OrderId:    msg.OrderId,
			FromType:   msg.FromType,
			ToType:     msg.FromType,
			MsgType:    msg.MsgType,
			CreateTime: msg.CreateTime,
			IsSystem:   MSG_SYSTEM,
			SystemType: MSG_SYSTEM_TYPE_ERROR,
			Source:     MSG_SOURCE_WX,
		}
		var strreturnmsg []byte
		strreturnmsg, _ = json.Marshal(returnmsg)
		h.message <- string(strreturnmsg)
	}

	return posterr
}

type PushResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func (m *Message) sendUserIOS() error {
	data := make(ur.Values)
	data.Set("deviceid", m.To)
	data.Set("content", m.Message)
	res, err := http.PostForm(fmt.Sprintf("%spush/ios/single", HTTP_SERVICE_URL), data)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var pr PushResponse
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return err
	}
	if !pr {
		return errors.New(pr.Message)
	}
	return nil
}

func (m *Message) sendUserAndroid() error {
	data := make(ur.Values)
	data.Set("deviceid", m.To)
	data.Set("content", m.Message)
	res, err := http.PostForm(fmt.Sprintf("%spush/android/single", HTTP_SERVICE_URL), data)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var pr PushResponse
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return err
	}
	if !pr {
		return errors.New(pr.Message)
	}
	return nil

}

func (m *Message) sendUser360Stream() error {
	return nil
}
