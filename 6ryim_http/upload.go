package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyf/postwx"
)

func uploadDir(prepath string) (string, error) {
	now := time.Now()
	path := fmt.Sprintf("%s/%v/%s/%v", prepath, now.Year(), now.Month(), now.Day())
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, nil
}

func init() {
	handlers["/upload"] = func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
		var result string

		err := r.ParseMultipartForm(C.maxImageSize)
		if err != nil {
			result = fmt.Sprintf("%v", err)
			response(w, result)
			return
		}

		if r.MultipartForm != nil && len(r.MultipartForm.File) > 0 {
			for _, files := range r.MultipartForm.File {
				if len(files) == 0 {
					result = "no file upload"
					response(w, result)
					return
				}
				file := files[0]
				f, err := file.Open()
				if err != nil {
					logger.Printf("upload file err:%v", err)
					result = "file error"
					response(w, result)
					return
				}

				d, err := ioutil.ReadAll(f)
				if err != nil {
					logger.Printf("upload file err:%v", err)
					result = "read file error"
					response(w, result)
					return
				}

				data_len := int64(len(d))
				if data_len > C.maxImageSize {
					result = "read file size more than maxImageSize"
					response(w, result)
					return
				}

				dir, err := uploadDir(C.uploadpath)
				if err != nil {
					logger.Printf("uploadDir err:%v", err)
					result = "upload dir error"
					response(w, result)
					return
				}

				if state := IsImage(file.Filename); !state {
					result = "upload file extension invalid"
					response(w, result)
					return
				}

				fp := fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%v%s", time.Now().UnixNano(), filepath.Ext(file.Filename)))
				ioutil.WriteFile(fp, d, os.ModePerm)
				result1 := map[string]string{
					"filepath": strings.Replace(fp, C.uploadpath, "", -1),
				}

				source := r.Form.Get("source")
				if strings.EqualFold(source, "wx") {
					mediaType := "image"
					media_id, err := postwx.UploadMedia(fp, mediaType)
					if err == nil {
						logger.Printf("postwx.UploadMedia err:%v", err)
						result1["media_id"] = media_id
					}
				}

				response(w, result1)
				break
			}
		} else {
			result = "no file upload"
			response(w, result)
			return
		}
	}

	handlers["/uploadwx"] = func(w http.ResponseWriter, r *http.Request, logger *log.Logger, params url.Values) {
		filepath := params.Get("filepath")
		mediaType := params.Get("mediatype")
		media_id, err := postwx.UploadMedia(filepath, mediaType)
		if err != nil {
			logger.Printf("postwx.UploadMedia err:%v", err)
			response(w, "Server Invalid")
			return
		}
		response(w, map[string]string{media_id: media_id})
	}
}
