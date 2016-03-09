package main

import (
	"fmt"
	"time"
)

func getFormatNow(tpl string) string {
	tpls := map[string]string{
		"zh":  "%v年%d月%v %v:%v:%v",
		"num": "%v-%d-%v %v:%v:%v",
	}
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	date := fmt.Sprintf(tpls[tpl], year, month, day, hour, minute, second)
	return date
}
