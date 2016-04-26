package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	start, end := 0, 69
	target := "http://77flz7.com1.z0.glb.clouddn.com/weixin-emotions/%d.png"

	for ; start <= end; start++ {
		res, err := http.Get(fmt.Sprintf(target, start))
		if err != nil {
			fmt.Printf("exit err:%v\n", err)
			continue
		}
		body, _ := ioutil.ReadAll(res.Body)
		err = ioutil.WriteFile(fmt.Sprintf("./%d.png", start), body, 0666)
		if err != nil {
			fmt.Printf("exit err:%v\n", err)
		}
		res.Body.Close()
	}
}
