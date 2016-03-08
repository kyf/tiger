package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-mgo/mgo/bson"
	im_type "github.com/kyf/6ryim/6ryim_http/im_type"
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
		panic(err)
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
}
