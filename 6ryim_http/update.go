package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-mgo/mgo"
	"github.com/go-mgo/mgo/bson"
)

type DeviceToken struct {
	Id             bson.ObjectId `bson:"_id" json:"id"`
	DeviceTokenStr string        `bson:"devicetoken" json:"devicetoken"`
	UserId         string        `bson:"userid" json:"userid"`
}

func init() {
	var document string = "devicetoken"

	handlers["/update"] = func(w http.ResponseWriter, r *http.Request, params url.Values, logger *log.Logger) {
		deviceToken := params.Get("deviceToken")
		userId := params.Get("userId")

		var result string
		if len(deviceToken) == 0 {
			result = "please pass the deviceToken"
			response(w, result)
			return
		}

		if len(userId) == 0 {
			result = "please pass the userId"
			response(w, result)
			return
		}

		cli := NewMongoClient()
		err := cli.Connect()
		if err != nil {
			logger.Printf("%v", err)
			result = "Server Invalid"
			response(w, result)
			return
		}
		defer cli.Close()

		err = cli.Remove(document, bson.M{"userid": userId})
		if err != nil {
			logger.Printf("%v", err)
		}

		bsonDeviceToken := &DeviceToken{
			bson.NewObjectId(),
			deviceToken,
			userId,
		}
		cli.Add(document, bsonDeviceToken)

		result1 := make(map[string]string)
		response(w, result1)

	}

	handlers["/getDeviceTokenByUserId"] = func(w http.ResponseWriter, r *http.Request, params url.Values, logger *log.Logger) {
		userId := params.Get("userId")

		var msg string
		if len(userId) == 0 {
			msg = "please pass the userId"
			response(w, msg)
			return
		}

		var result []DeviceToken

		cli := NewMongoClient()
		err := cli.Connect()
		if err != nil {
			logger.Printf("%v", err)
			msg = "Server Invalid"
			response(w, result)
			return
		}
		defer cli.Close()
		err = cli.Find(document, bson.M{"userid": userId}).Limit(1).All(&result)
		if err != nil && err != mgo.ErrNotFound {
			logger.Printf("%v", err)
			msg = "Server Invalid"
			response(w, result)
			return
		}

		if len(result) > 0 {
			result1 := map[string]string{
				"deviceToken": result[0].DeviceTokenStr,
			}
			response(w, result1)
		} else {
			response(w, "deviceToken not found!")
		}
	}
}
