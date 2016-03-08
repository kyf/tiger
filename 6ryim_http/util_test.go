package main

import (
	"fmt"
	"testing"
)

func TestPathinfo(t *testing.T) {
	var re map[string]string
	var err error

	re, err = Pathinfo("")
	fmt.Println(re, err)

	re, err = Pathinfo("asdasd")
	fmt.Println(re, err)

	re, err = Pathinfo("/home/user/kyf")
	fmt.Println(re, err)

	re, err = Pathinfo("/usr/local/kyf.jpeg")
	fmt.Println(re, err)
}

func TestIsImage(t *testing.T) {
	var re bool

	re = IsImage("/usr/local/kyf.js")
	fmt.Println(re)

	re = IsImage("/usr/local/kyf.gif")
	fmt.Println(re)

	re = IsImage("/usr/local/kyf.JPEG")
	fmt.Println(re)
}
