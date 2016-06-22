package main

import (
	"github.com/Unknwon/goconfig"
)

type Config struct {
	uploadpath   string
	port         string
	sslport      string
	maxImageSize int64

	redisServer     string
	redisPass       string
	redisMaxIdleNum int

	mongodbServer   string
	mongodbPort     string
	mongodbName     string
	msgMongodbName  string
	mongodbUser     string
	mongodbPass     string
	mongodbPoolSize int
}

func initConfig(fileName string) (Config, error) {
	var conf Config
	c, err := goconfig.LoadConfigFile(fileName)
	if err != nil {
		return conf, err
	}

	conf.uploadpath, err = c.GetValue("default", "uploadpath")
	if err != nil {
		return conf, err
	}
	conf.port = c.MustValue("default", "port", "8989")
	conf.sslport = c.MustValue("default", "sslport", "4443")
	conf.maxImageSize = c.MustInt64("default", "max_image_size", 1024*1024*2)

	conf.redisServer = c.MustValue("redis", "redis_server", "")
	conf.redisPass = c.MustValue("redis", "redis_pass", "")
	conf.redisMaxIdleNum = c.MustInt("redis", "redis_max_idle_num", 10)

	conf.mongodbServer = c.MustValue("mongodb", "mongodb_server", "")
	conf.mongodbPort = c.MustValue("mongodb", "mongodb_port", "")
	conf.mongodbUser = c.MustValue("mongodb", "mongodb_user", "")
	conf.mongodbName = c.MustValue("mongodb", "mongodb_name", "")
	conf.msgMongodbName = c.MustValue("mongodb", "msgmongodb_name", "")
	conf.mongodbPass = c.MustValue("mongodb", "mongodb_pass", "")
	conf.mongodbPoolSize = c.MustInt("mongodb", "mongodb_pool_size", 300)

	InitRedis(conf.redisServer, conf.redisPass, conf.redisMaxIdleNum)

	InitMongodb(conf.mongodbServer, conf.mongodbPort, conf.mongodbUser, conf.mongodbName, conf.msgMongodbName, conf.mongodbPass, conf.mongodbPoolSize)
	return conf, nil
}
