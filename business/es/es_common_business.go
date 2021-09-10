package es

import (
	"dongchamao/global/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

func GetESTableByTime(table string, start, stop time.Time) (string, error) {
	//时间截止至9.1号
	if start.Unix() < 1630425600 {
		start = time.Unix(1630425600, 0)
	}
	if start.After(stop) {
		return "", errors.New("参数错误")
	}
	if start.Format("20060102") == stop.Format("20060102") {
		return fmt.Sprintf(table, start.Format("20060102")), nil
	}
	esTableArr := make([]string, 0)
	begin, _ := time.ParseInLocation("20060102", start.Format("200601")+"01", time.Local)
	endMonth := stop.Month()
	year := utils.ToInt(stop.Format("2006")) - utils.ToInt(start.Format("2006"))
	endMonth += 12 * time.Month(year)
	beginMonth := begin.Month()
	for {
		if beginMonth > endMonth {
			break
		}
		esTableArr = append(esTableArr, fmt.Sprintf(table, begin.Format("200601")+"*"))
		begin = begin.AddDate(0, 1, 0)
		beginMonth += 1
	}
	return strings.Join(esTableArr, ","), nil
}

func GetESTableByDayTime(table string, start, stop time.Time) (string, error) {
	//时间截止至9.1号
	if start.Unix() < 1630425600 {
		start = time.Unix(1630425600, 0)
	}
	if start.After(stop) {
		return "", errors.New("参数错误")
	}
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
	return strings.Join(esTableArr, ","), nil
}

func GetESTableByMonthTime(table string, start, stop time.Time) (string, error) {
	//时间截止至9.1号
	if start.Unix() < 1630425600 {
		start = time.Unix(1630425600, 0)
	}
	if start.After(stop) {
		return "", errors.New("参数错误")
	}
	esTableArr := make([]string, 0)
	begin, _ := time.ParseInLocation("20060102", start.Format("200601")+"01", time.Local)
	endMonth := stop.Month()
	year := utils.ToInt(stop.Format("2006")) - utils.ToInt(start.Format("2006"))
	endMonth += 12 * time.Month(year)
	beginMonth := begin.Month()
	for {
		if beginMonth > endMonth {
			break
		}
		esTableArr = append(esTableArr, fmt.Sprintf(table, begin.Format("200601")))
		beginMonth += 1
		begin = begin.AddDate(0, 1, 0)
	}
	return strings.Join(esTableArr, ","), nil
}
