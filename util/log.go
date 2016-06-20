package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	LAYOUT = "2006-01-02"
)

func getTodayLogName(prefix string) string {
	return prefix + time.Now().Format(LAYOUT) + ".log"
}

type Writer struct {
	name   string
	prefix string
	fp     *os.File
	mutex  sync.Mutex
}

func (w *Writer) Write(data []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	today := getTodayLogName(w.prefix)
	if !strings.EqualFold(today, w.name) {
		var fp *os.File
		fp, err = os.OpenFile(today, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return
		}
		w.fp.Close()
		w.name = today
		w.fp = fp
	}

	n, err = w.fp.Write(data)
	return
}

func (w *Writer) Close() (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	err = w.fp.Close()
	return
}

func NewWriter(name string) (io.WriteCloser, error) {
	prefix := filepath.Dir(name) + "/"
	today := getTodayLogName(prefix)
	fp, err := os.OpenFile(today, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Writer{name: today, prefix: prefix, fp: fp}, nil
}
