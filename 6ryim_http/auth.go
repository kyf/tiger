package main

import (
	"log"
	"net/http"
	"net/url"
)

func init() {
	handlers["/auth"] = func(w http.ResponseWriter, r *http.Request, params url.Values, logger *log.Logger) {

		var result interface{}
		clientid := params.Get("deviceid")
		systemVersion := params.Get("systemVersion")
		systemName := params.Get("systemName")
		deviceModel := params.Get("deviceModel")
		country := params.Get("country")
		language := params.Get("language")
		timezone := params.Get("timezone")
		name := params.Get("name")
		appVersion := params.Get("appVersion")

		if len(clientid) == 0 {
			result = "please pass the deviceid"
			response(w, result)
			return
		}

		if stoken, err := getTokenByDeviceId(clientid); err == nil {
			result = map[string]string{
				"token": stoken,
			}
			response(w, result)
			return
		}

		token := newToken(clientid, name, systemVersion, systemName, deviceModel, country, language, timezone, appVersion)
		err := token.Connect()
		if err != nil {
			logger.Printf("token connect err:%v", err)
			result = "server error"
			response(w, result)
			return
		}

		defer token.Disconnect()
		err = token.Fresh()
		if err != nil {
			logger.Printf("token fresh err:%v", err)
			result = "server error"
			response(w, result)
			return
		}

		result = map[string]string{
			"token": token.GetToken(),
		}
		response(w, result)
	}
}
