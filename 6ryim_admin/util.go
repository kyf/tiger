package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

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
