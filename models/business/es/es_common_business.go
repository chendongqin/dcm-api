package es

import (
	"fmt"
	"strings"
	"time"
)

func GetESTableByTime(table string, start, stop time.Time) string {
	esTableArr := make([]string, 0)
	begin := start
	endMonth := stop.Month()
	for {
		if begin.Month() > endMonth {
			break
		}
		esTableArr = append(esTableArr, fmt.Sprintf(table, begin.Format("200601")+"*"))
		begin = begin.AddDate(0, 1, 0)
	}
	return strings.Join(esTableArr, ",")
}
