package util

import (
	"sync"
	"testing"
)

func TestWrite(t *testing.T) {
	p := "/var/log/6ryim_test/6ryim_test.log"
	w, err := NewWriter(p)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	var wg sync.WaitGroup
	for i := 0; i <= 500; i++ {
		go func() {
			for k := 0; k < 10000; k++ {
				_, err = w.Write([]byte("this is a log\n"))
				if err != nil {
					t.Error(err)
				}
			}
			wg.Done()
		}()
		wg.Add(1)
	}

	wg.Wait()
}
