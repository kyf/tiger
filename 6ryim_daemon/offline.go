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

func fetchOffline(msg Message) ([]string, error) {

}

func countOffline(msg Message) (int, error) {

}
