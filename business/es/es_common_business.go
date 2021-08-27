package es

import (
	"fmt"
	"strings"
	"time"
)

func GetESTableByTime(table string, start, stop time.Time) string {
	esTableArr := make([]string, 0)
	begin := start
	endTime := stop.Unix()
	for {
		if begin.Unix() > endTime {
			break
		}
		esTableArr = append(esTableArr, fmt.Sprintf(table, begin.Format("200601")+"*"))
		begin = begin.AddDate(0, 1, 0)
	}
	return strings.Join(esTableArr, ",")
}

func GetESTableByDayTime(table string, start, stop time.Time) string {
	esTableArr := make([]string, 0)
	begin := start
	endDay := stop.Unix()
	for {
		if begin.Unix() > endDay {
			break
		}
		esTableArr = append(esTableArr, fmt.Sprintf(table, begin.Format("20060102")))
		begin = begin.AddDate(0, 0, 1)
	}
	return strings.Join(esTableArr, ",")
}
