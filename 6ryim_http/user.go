package main

import (
	"fmt"
	"net/http"
	"net/url"

	MyRedis "github.com/garyburd/redigo/redis"
)

func init() {
	handlers["/user"] = func(w http.ResponseWriter, r *http.Request, params url.Values) {
		var result string
		token := params.Get("token")
		if len(token) == 0 {
			result = "please pass the token"
			response(w, result)
			return
		}

		redis, err := GetRedis()
		if err != nil {
			result = "server error"
			response(w, result)
			return
		}

		defer redis.Close()
		value, err := MyRedis.StringMap(redis.Do("hgetall", fmt.Sprintf("%s%s", GetSessionPrefix(), token)))
		if err != nil || len(value) == 0 {
			result = "token invalid"
			response(w, result)
			return
		}
		response(w, value)
	}
}
