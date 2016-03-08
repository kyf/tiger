package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	handlers[".*(jpg|gif|png|bmp|amr)"] = func(w http.ResponseWriter, r *http.Request) {
		path := fmt.Sprintf("%s/%s", C.uploadpath, r.URL.Path)
		file, err := os.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		buffer := bytes.NewBuffer(nil)
		_, err = io.Copy(buffer, file)
		if err != nil {
			response(w, "Server Error")
			return
		}

		ext := filepath.Ext(path)

		w.Header().Add("Content-Type", mediatype[ext])
		buffer.WriteTo(w)
	}
}
