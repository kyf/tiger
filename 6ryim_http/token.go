package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	SESSION_PREFIX      string = "session_"
	OFFLINE_MSG_PREFIX  string = "offline_msg_"
	DEVICE_TOKEN_PREFIX string = "device_token_"
)

type token struct {
	clientid      string
	value         string
	conn          redis.Conn
	Name          string
	SystemVersion string
	SystemName    string
	DeviceModel   string
	Country       string
	Language      string
	TimeZone      string
	AppVersion    string
}

func getTokenByDeviceId(deviceid string) (token string, err error) {
	redis_cli, err := GetRedis()
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s%s", DEVICE_TOKEN_PREFIX, deviceid)
	token, err = redis.String(redis_cli.Do("get", key))
	if err != nil {
		return "", err
	}
	return token, nil
}

func newToken(clientid, name, systemVersion, systemName, deviceModel, country, language, timezone, appVersion string) *token {
	timestamp := time.Now().UnixNano()
	d := fmt.Sprintf("%v%s", timestamp, clientid)
	v := md5.Sum([]byte(d))
	va := hex.EncodeToString(v[0:])
	return &token{
		clientid:      clientid,
		value:         va,
		Name:          name,
		SystemVersion: systemVersion,
		SystemName:    systemName,
		DeviceModel:   deviceModel,
		Country:       country,
		Language:      language,
		TimeZone:      timezone,
		AppVersion:    appVersion,
	}
}

type Redis struct {
	uri  string
	pass string
}

var (
	rediss    Redis
	redisPool *redis.Pool
)

func InitRedis(redisServer, redisPass string, redisMaxIdleNum int) {
	rediss = Redis{redisServer, redisPass}
	redisPool = &redis.Pool{
		MaxIdle:     redisMaxIdleNum,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rediss.uri)
			if err != nil {
				return nil, err
			}

			_, err = c.Do("auth", rediss.pass)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

}

func GetRedis() (redis.Conn, error) {
	return redisPool.Get(), nil
}

func GetSessionPrefix() string {
	return SESSION_PREFIX
}

func (t *token) Connect() error {
	c, err := GetRedis()
	if err != nil {
		return err
	}
	t.conn = c
	return nil
}

func (t *token) Fresh() error {
	cmd := "hset"
	key := fmt.Sprintf("%s%s", SESSION_PREFIX, t.value)
	args := map[string]string{
		"clientid":      t.clientid,
		"value":         t.value,
		"name":          t.Name,
		"systemversion": t.SystemVersion,
		"systemname":    t.SystemName,
		"devicemodel":   t.DeviceModel,
		"country":       t.Country,
		"language":      t.Language,
		"timezone":      t.TimeZone,
		"appversion":    t.AppVersion,
	}
	for k, v := range args {
		_, err := t.conn.Do(cmd, key, k, v)
		if err != nil {
			return err
		}
	}

	arg := 60 * 60 * 24
	expire(key, arg, t.conn)

	key = fmt.Sprintf("%s%s", DEVICE_TOKEN_PREFIX, t.clientid)
	_, err := redis.String(t.conn.Do("set", key, t.value))
	if err != nil {
		return err
	}
	expire(key, arg, t.conn)
	return nil
}

func expire(key string, arg int, conn redis.Conn) error {
	cmd := "expire"
	_, err := conn.Do(cmd, key, arg)
	if err != nil {
		return err
	}
	return nil
}

func (t *token) GetToken() string {
	return t.value
}

func (t *token) GetClientId() string {
	return t.clientid
}

func (t *token) Del() {
	cmd := "del"
	key := fmt.Sprintf("%s%s", SESSION_PREFIX, t.value)
	t.conn.Do(cmd, key)
}

func (t *token) Disconnect() {
	t.conn.Close()
}
