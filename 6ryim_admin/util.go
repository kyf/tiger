package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kyf/postwx"
)

func fetchWxMedia(fpath string) (string, error) {
	data := make(url.Values)
	data.Set("fpath", fpath)
	dir, err := uploadDir(UPLOAD_PATH)
	if err != nil {
		return "", err
	}

	fp := fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%v", time.Now().UnixNano()))
	fullpath, err := postwx.GetMedia(fpath, fp)
	if err != nil {
		return "", err
	}
	newpath := strings.Replace(string(fullpath), UPLOAD_PATH, "", -1)
	return newpath, nil
}

func Pathinfo(file string) (map[string]string, error) {
	if len(file) == 0 {
		return nil, errors.New("file path is empty!")
	}
	var result map[string]string = make(map[string]string)
	info := strings.Split(file, "/")
	info_len := len(info)
	if info_len == 1 {
		result["basepath"] = ""
		result["filename"] = file
	} else {
		result["filename"] = info[info_len-1]
		result["basepath"] = strings.Replace(file, result["filename"], "", -1)
	}
	name := strings.Split(result["filename"], ".")
	result["extension"] = ""
	name_len := len(name)
	if name_len > 1 {
		result["extension"] = name[name_len-1]
	}

	return result, nil
}

func IsImage(file string) bool {
	pathinfo, err := Pathinfo(file)
	if err != nil {
		return false
	}
	exts := []string{"jpeg", "jpg", "png", "gif"}
	extension := strings.ToLower(pathinfo["extension"])
	for _, e := range exts {
		if strings.EqualFold(extension, e) {
			return true
		}
	}
	return false
}

func cbResponseJson(writer io.Writer, status bool, msg, cb string, data ...interface{}) {
	if len(msg) == 0 {
		msg = "success"
	}
	result := map[string]interface{}{
		"status": status,
		"msg":    msg,
	}

	if len(data) > 0 {
		result["data"] = data[0]
	}

	re, _ := json.Marshal(result)
	writer.Write([]byte(fmt.Sprintf("<script type='text/javascript'>%s(%s)</script>", cb, string(re))))
}

func responseJson(writer io.Writer, status bool, msg string, data ...interface{}) {
	if len(msg) == 0 {
		msg = "success"
	}
	result := map[string]interface{}{
		"status": status,
		"msg":    msg,
	}

	if len(data) > 0 {
		result["data"] = data[0]
	}

	re, _ := json.Marshal(result)
	writer.Write(re)
}

func fetchTop(admin_name string) ([]byte, error) {
	top, err := fetchFile("./tpl/top.html")
	if err != nil {
		return nil, err
	}
	_top := strings.Replace(string(top), "{admin_name}", admin_name, -1)
	return []byte(_top), nil
}

func fetchLeft() ([]byte, error) {
	return fetchFile("./tpl/left.html")
}

func fetchFile(path string) ([]byte, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	content, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func StringSliceContains(it string, its []string) bool {
	for _, item := range its {
		if strings.EqualFold(item, it) {
			return true
		}
	}

	return false
}

func getAccessToken() (string, error) {
	res, err := http.Get("http://m.6renyou.com/weixin_service/getAccessToken?account_type=1")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
