package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	DEFAULT_POOL_EXPIRES time.Duration = time.Second * 60 * 60 * 2

	USER_REQUEST_URL string = "http://admin.6renyou.com/weixin_plugin/getAjaxMemberInfo"
)

type User struct {
	RealName string `json:"realname"`
	Mobile   string `json:"mobile"`
	WxOpenid string `json:"wx_openid"`
	UserID   string `json:"userid"`

	created time.Time
}

type UserManager struct {
	userPool map[string]User
	expires  time.Duration
	mutex    sync.Mutex
}

func NewUserManager() *UserManager {
	return &UserManager{userPool: make(map[string]User), expires: DEFAULT_POOL_EXPIRES}
}

func (um *UserManager) Get(userids, source []string) ([]User, error) {
	result := make([]User, 0, len(userids))
	bad := make(map[string][]string)
	bad["app"] = make([]string, 0, 4)
	bad["weixin"] = make([]string, 0, 4)
	for index, userid := range userids {
		if u, ok := um.userPool[userid]; ok {
			if u.created.Add(um.expires).Unix() > time.Now().Unix() {
				result = append(result, u)
			} else {
				um.mutex.Lock()
				delete(um.userPool, userid)
				um.mutex.Unlock()
				bad[source[index]] = append(bad[source[index]], userid)
			}
		} else {
			bad[source[index]] = append(bad[source[index]], userid)
		}
	}

	if len(bad["app"]) == 0 && len(bad["weixin"]) == 0 {
		return result, nil
	}

	badResult, err := requestUser(bad)
	if err != nil {
		return nil, err
	}
	um.mutex.Lock()
	for _, u := range badResult {
		u.created = time.Now()
		um.userPool[u.UserID] = u
		result = append(result, u)
	}
	um.mutex.Unlock()

	return result, nil
}

func requestUser(userids map[string][]string) ([]User, error) {
	params := make(url.Values)
	params.Set("ids", strings.Join(userids["app"], ","))
	params.Set("wxids", strings.Join(userids["weixin"], ","))
	res, err := http.PostForm(USER_REQUEST_URL, params)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := struct {
		m6ryResponse
		Data map[string][]User `json:"data"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != 1 {
		return nil, errors.New(result.Info)
	}

	appUsers := result.Data["app"]
	weixinUsers := result.Data["weixin"]

	data := make([]User, len(appUsers)+len(weixinUsers))
	var n int
	n = copy(data[n:], appUsers)
	copy(data[n:], weixinUsers)
	return data, nil
}
