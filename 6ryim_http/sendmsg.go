package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	im_type "github.com/kyf/6ryim/6ryim_http/im_type"
)

func init() {
	handlers["/sendmsg"] = func(w http.ResponseWriter, r *http.Request, params url.Values, logger *log.Logger) {
		msg := params.Get("msg")
		var result string
		if len(msg) == 0 {
			result = "msg is empty!"
			response(w, result)
			return
		}

		go processMessage(msg, logger)
		result1 := make(map[string]string)
		result1["response"] = "processMessage"
		response(w, result1)
	}

}

func processMessage(msg string, logger *log.Logger) {
	var message im_type.Message
	err := json.Unmarshal([]byte(msg), &message)
	if err != nil {
		logger.Printf("processMessage json.Unmarshal err:%v", err)
		return
	}

	var param url.Values = make(url.Values)
	param.Set("msg", msg)
	res, err := http.PostForm(fmt.Sprintf("%smessage/receive", WS_SERVICE_URL), param)
	if err != nil {
		logger.Printf("processMessage.postform err:%v", err)
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Printf("processMessage readall err:%v", err)
		return
	}
	logger.Printf("[processMessage]response is %s", string(data))
}
