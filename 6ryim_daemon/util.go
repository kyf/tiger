package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func getFormatNow(tpl string) string {
	tpls := map[string]string{
		"zh":  "%v年%d月%v %v:%v:%v",
		"num": "%v-%d-%v %v:%v:%v",
	}
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	date := fmt.Sprintf(tpls[tpl], year, month, day, hour, minute, second)
	return date
}

func response(w http.ResponseWriter, status bool, msg string, data ...interface{}) {
	result := map[string]interface{}{
		"status": status,
		"msg":    msg,
	}

	if len(data) > 0 {
		result["data"] = data[0]
	}

	res, err := json.Marshal(result)
	if err != nil {
		w.Write([]byte("Server Invalid"))
		return
	}

	w.Write(res)
}

func generalSign(accessid, timestamp, secretkey string) string {
	origin := fmt.Sprintf("%s%s%s", timestamp, accessid, secretkey)
	tmp := md5.Sum([]byte(origin))
	return hex.EncodeToString(tmp[:])
}
