package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	im_type "github.com/kyf/6ryim/6ryim_http/im_type"
	"gopkg.in/mgo.v2/bson"
)

var (
	cli Mongo
)

func getData(params url.Values) ([]im_type.BsonMessage, error) {
	lastid := params.Get("lastid")
	_size := params.Get("size")
	key := params.Get("key")
	from := params.Get("from")
	to := params.Get("to")
	orderid := params.Get("orderid")
	sort := params.Get("sort")

	cli := NewMongoClient()
	err := cli.Connect()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	size := 20
	if len(_size) > 0 {
		size, err = strconv.Atoi(_size)
		if err != nil {
			size = 10
		}
	}

	var result []im_type.BsonMessage
	where := make([]bson.M, 0)

	if len(lastid) > 0 {
		where = append(where, bson.M{"_id": bson.M{"$lt": bson.ObjectIdHex(lastid)}})
	}

	if len(key) > 0 {
		where = append(where, bson.M{"message": bson.M{"$regex": key}})
	}

	if len(from) > 0 {
		where = append(where, bson.M{"from": from})
	}

	if len(to) > 0 {
		where = append(where, bson.M{"to": to})
	}

	if len(orderid) > 0 {
		where = append(where, bson.M{"orderid": orderid})
	}

	var condition bson.M = nil
	if len(where) > 0 {
		if len(where) > 1 {
			condition = bson.M{"$and": where}
		} else {
			condition = where[0]
		}
	}

	var sort_cond string = "-_id"
	if sortd, err := strconv.Atoi(sort); err == nil {
		if sortd > 0 {
			sort_cond = "_id"
		}
	}

	err = cli.Find("message", condition).Sort(sort_cond).Limit(size).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getNewMessageNum(params url.Values) (int, error) {
	lastid := params.Get("lastid")
	orderid := params.Get("orderid")
	fromtype := params.Get("fromtype")

	cli := NewMongoClient()
	err := cli.Connect()
	if err != nil {
		return 0, err
	}
	defer cli.Close()

	where := make([]bson.M, 0)

	if len(lastid) > 0 {
		where = append(where, bson.M{"_id": bson.M{"$gt": bson.ObjectIdHex(lastid)}})
	}

	if len(orderid) > 0 {
		where = append(where, bson.M{"orderid": orderid})
	}

	if len(fromtype) > 0 {
		where = append(where, bson.M{"fromtype": fromtype})
	}

	var condition bson.M = nil
	if len(where) > 0 {
		if len(where) > 1 {
			condition = bson.M{"$and": where}
		} else {
			condition = where[0]
		}
	}

	var total int
	total, err = cli.Find("message", condition).Count()
	if err != nil {
		return 0, err
	}

	return total, nil
}

func getPageData(params url.Values) ([]im_type.BsonMessage, int, error) {
	_size := params.Get("size")
	_page := params.Get("page")
	key := params.Get("key")
	from := params.Get("from")
	fromtype := params.Get("fromtype")
	to := params.Get("to")
	orderid := params.Get("orderid")
	msgtype := params.Get("msgtype")
	msgsource := params.Get("msgsource")
	sort := params.Get("sort")

	cli := NewMongoClient()
	err := cli.Connect()
	if err != nil {
		return nil, 0, err
	}
	defer cli.Close()

	size := 20
	if len(_size) > 0 {
		size, err = strconv.Atoi(_size)
		if err != nil {
			size = 20
		}
	}

	page := 1
	if len(_page) > 0 {
		page, err = strconv.Atoi(_page)
		if err != nil {
			page = 1
		}
	}

	var result []im_type.BsonMessage
	where := make([]bson.M, 0)

	if len(key) > 0 {
		where = append(where, bson.M{"message": bson.M{"$regex": key}})
	}

	if len(from) > 0 {
		where = append(where, bson.M{"from": from})
	}

	if len(fromtype) > 0 {
		where = append(where, bson.M{"fromtype": fromtype})
	}

	if len(to) > 0 {
		where = append(where, bson.M{"to": to})
	}

	if len(orderid) > 0 {
		where = append(where, bson.M{"orderid": orderid})
	}

	if len(msgtype) > 0 {
		where = append(where, bson.M{"msgtype": msgtype})
	}

	if len(msgsource) > 0 {
		where = append(where, bson.M{"source": msgsource})
	}

	var condition bson.M = nil
	if len(where) > 0 {
		if len(where) > 1 {
			condition = bson.M{"$and": where}
		} else {
			condition = where[0]
		}
	}

	var sort_cond string = "-_id"
	if sortd, err := strconv.Atoi(sort); err == nil {
		if sortd > 0 {
			sort_cond = "_id"
		}
	}

	skip := (page - 1) * size

	err = cli.Find("message", condition).Sort(sort_cond).Skip(skip).Limit(size).All(&result)
	if err != nil {
		return nil, 0, err
	}

	var total int
	total, err = cli.Find("message", condition).Sort(sort_cond).Count()
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func init() {
	handlers["/message"] = func(w http.ResponseWriter, logger *log.Logger, r *http.Request, params url.Values) {
		data, err := getData(params)
		if err != nil {
			logger.Printf("message.getData err:%v", err)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		} else {
			js, err := json.Marshal(data)
			if err == nil {
				w.Write(js)
			} else {
				w.Write([]byte(fmt.Sprintf("%v", err)))
			}
		}
	}

	handlers["/message/show"] = func(w http.ResponseWriter, logger *log.Logger, r *http.Request, params url.Values) {
		data, total, err := getPageData(params)
		if err != nil {
			logger.Printf("message.getPageData err:%v", err)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		} else {
			js, err := json.Marshal(map[string]interface{}{"data": data, "total": total})
			if err == nil {
				w.Write(js)
			} else {
				w.Write([]byte(fmt.Sprintf("%v", err)))
			}
		}
	}

	handlers["/message/new/number"] = func(w http.ResponseWriter, logger *log.Logger, r *http.Request, params url.Values) {
		data, err := getNewMessageNum(params)
		if err != nil {
			logger.Printf("message.getNewMessageNum err:%v", err)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		} else {
			js, err := json.Marshal(map[string]interface{}{"data": data})
			if err == nil {
				w.Write(js)
			} else {
				w.Write([]byte(fmt.Sprintf("%v", err)))
			}
		}
	}

	handlers["/store"] = func(params url.Values, w http.ResponseWriter, r *http.Request, logger *log.Logger) {
		m := params.Get("msg")

		if strings.EqualFold("", m) {
			response(w, "msg is empty")
			return
		}

		var msg im_type.Message
		err := json.Unmarshal([]byte(m), &msg)
		if err != nil {
			logger.Printf("msg [%s] json.Unmarshal err:%v", m, err)
			response(w, "msg Invalid")
			return
		}

		mongocli := NewMongoClient()
		mongocli.Connect()
		defer mongocli.Close()

		var bmsg *im_type.BsonMessage = &im_type.BsonMessage{
			bson.NewObjectId(),
			msg.From,
			msg.To,
			msg.Message,
			msg.OrderId,
			msg.FromType,
			msg.ToType,
			msg.MsgType,
			msg.CreateTime,
			msg.IsSystem,
			msg.SystemType,
			msg.Source,
		}

		err = mongocli.Add("message", bmsg)
		if err != nil {
			logger.Printf("store message err:%v", err)
			response(w, "Server Invalid")
			return
		}
		response(w, make(map[string]string))
	}

	handlers["/offline/store"] = func(params url.Values, w http.ResponseWriter, r *http.Request, logger *log.Logger) {
		msg := params.Get("msg")
		to := params.Get("to")

		if strings.EqualFold("", msg) {
			response(w, "msg is empty")
			return
		}

		if strings.EqualFold("", to) {
			response(w, "to is empty")
			return
		}

		redis_cli, err := GetRedis()
		if err != nil {
			logger.Printf("redis connection error :%v", err)
			response(w, "Server Invalid")
			return
		}
		key := fmt.Sprintf("%s%s", OFFLINE_MSG_PREFIX, to)
		_, err = redis_cli.Do("rpush", key, msg)
		if err != nil {
			logger.Printf("redis store offline message  error :%v", err)
			response(w, "Server Invalid")
			return
		}

		response(w, make(map[string]string))
	}

	handlers["/offline/fetch"] = func(params url.Values, w http.ResponseWriter, r *http.Request, logger *log.Logger) {
		to := params.Get("to")

		if strings.EqualFold("", to) {
			response(w, "to is empty")
			return
		}

		redis_cli, err := GetRedis()
		if err != nil {
			logger.Printf("redis connection error :%v", err)
			response(w, "Server Invalid")
			return
		}
		key := fmt.Sprintf("%s%s", OFFLINE_MSG_PREFIX, to)
		v, err := redis.Strings(redis_cli.Do("lrange", key, 0, -1))
		if err != nil {
			logger.Printf("redis fetch offline message  error :%v", err)
			response(w, "Server Invalid")
			return
		}
		redis_cli.Do("del", key)

		response(w, map[string][]string{"data": v})
	}

	handlers["/offline/count"] = func(params url.Values, w http.ResponseWriter, r *http.Request, logger *log.Logger) {
		to := params.Get("to")

		if strings.EqualFold("", to) {
			response(w, "to is empty")
			return
		}

		redis_cli, err := GetRedis()
		if err != nil {
			logger.Printf("redis connection error :%v", err)
			response(w, "Server Invalid")
			return
		}
		key := fmt.Sprintf("%s%s", OFFLINE_MSG_PREFIX, to)
		v, err := redis.Int(redis_cli.Do("llen", key))
		if err != nil {
			logger.Printf("redis count offline message  error :%v", err)
			response(w, "Server Invalid")
			return
		}

		response(w, map[string]int{"count": v})
	}
}
