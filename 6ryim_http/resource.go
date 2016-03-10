package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyf/postwx"
)

var (
	mediatype map[string]string = map[string]string{
		".jpg": "image/jpeg",
		".gif": "image/gif",
		".png": "image/png",
		".amr": "audio/amr",
	}
)

func init() {
	handlers[".*(jpg|gif|png|bmp|amr)"] = func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
		path := fmt.Sprintf("%s/%s", C.uploadpath, r.URL.Path)
		file, err := os.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		buffer := bytes.NewBuffer(nil)
		_, err = io.Copy(buffer, file)
		if err != nil {
			logger.Printf("read image file err:%v", err)
			response(w, "Server Error")
			return
		}

		ext := filepath.Ext(path)

		w.Header().Add("Content-Type", mediatype[ext])
		buffer.WriteTo(w)
	}

	handlers["/wxmedia/fetch"] = func(w http.ResponseWriter, r *http.Request, params url.Values, logger *log.Logger) {
		fpath := params.Get("fpath")

		dir, err := uploadDir(C.uploadpath)
		if err != nil {
			logger.Printf("uploadDir err:%v", err)
			response(w, "Server Invalid")
			return
		}

		fp := fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%v", time.Now().UnixNano()))
		fullpath, err := postwx.GetMedia(fpath, fp)
		if err != nil {
			logger.Printf("postwx.GetMedia error : %v", err)
			response(w, "Server Invalid")
			return
		}
		newpath := strings.Replace(string(fullpath), C.uploadpath, "", -1)
		response(w, map[string]string{"newpath": newpath})
	}
}
