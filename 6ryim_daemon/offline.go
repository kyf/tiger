package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func storeOffline(msg Message) error {
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data := make(url.Values)
	data.Set("msg", string(m))
	res, err = http.PostForm(fmt.Sprintf("%soffline/store", HTTP_SERVICE_URL), data)
	if err != nil {
		return err
	}
}

func fetchOffline(to string) ([]string, error) {
	data := make(url.Values)
	data.Set("to", to)
	res, err = http.PostForm(fmt.Sprintf("%soffline/fetch", HTTP_SERVICE_URL), data)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	type tdata struct {
		Status string              `json:"status"`
		Msg    string              `json:"msg"`
		Data   map[string][]string `json:"data"`
	}

	var result tdata
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold("ok", result.Status) {
		return nil, errors.New(result.Msg)
	}

	if list, ok := result.Data["data"]; ok {
		return list, nil
	} else {
		return nil, nil
	}

}

func countOffline(to string) (int, error) {
	data := make(url.Values)
	data.Set("to", to)
	res, err = http.PostForm(fmt.Sprintf("%soffline/count", HTTP_SERVICE_URL), data)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	type tdata struct {
		Status string         `json:"status"`
		Msg    string         `json:"msg"`
		Data   map[string]int `json:"data"`
	}

	var result tdata
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	if !strings.EqualFold("ok", result.Status) {
		return 0, errors.New(result.Msg)
	}

	if count, ok := result.Data["count"]; ok {
		return count, nil
	} else {
		return 0, nil
	}

}

func getDevicetokenByToken(token string) ([]byte, error) {
	data := make(url.Values)
	data.Set("token", token)
	res, err := http.PostForm(fmt.Sprintf("%sgetDevicetokenByToken", HTTP_SERVICE_URL), data)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if devicetoken, ok := result["devicetoken"]; ok {
		return []byte(devicetoken), nil
	} else {
		return nil, nil
	}
}
