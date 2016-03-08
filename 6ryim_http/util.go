package main

import (
	"errors"
	"strings"
)

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
