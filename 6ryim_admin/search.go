package main

import ()

func handleListAllMessage(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	msg_type, keyword := r.Form.Get("msg_type"), r.Form.Get("keyword")

	if len(openid) == 0 {
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
