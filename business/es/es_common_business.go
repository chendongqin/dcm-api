package es

import (
	"dongchamao/global/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

const DataStartTimestamp = 1627747200

var EsTableConnectionMap = map[string]string{
	"dy_product_aweme_author_analysis_%s": "aweme",
	"dy_product_aweme_%s":                 "aweme",
	"dy_aweme_%s":                         "aweme",
	"dy_aweme_product_%s":                 "aweme",
}

//对应的es集群
func SureConnection(tableName string) string {
	if v, ok := EsTableConnectionMap[tableName]; ok {
		return v
	}
	return ""
}

func GetESTableByTime(table string, start, stop time.Time) (string, string, error) {
	//时间截止至8.1号
	if start.Unix() < DataStartTimestamp {
		start = time.Unix(DataStartTimestamp, 0)
	}
	if start.After(stop) {
		return "", "", errors.New("参数错误")
	}
	if start.Format("20060102") == stop.Format("20060102") {
		return fmt.Sprintf(table, start.Format("20060102")), SureConnection(table), nil
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
	return strings.Join(esTableArr, ","), SureConnection(table), nil
}

func GetESTableByDayTime(table string, start, stop time.Time) (string, string, error) {
	//时间截止至8.1号
	if start.Unix() < DataStartTimestamp {
		start = time.Unix(DataStartTimestamp, 0)
	}
	if start.After(stop) {
		return "", "", errors.New("参数错误")
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
	return strings.Join(esTableArr, ","), SureConnection(table), nil
}

func GetESTableByMonthTime(table string, start, stop time.Time) (string, string, error) {
	//时间截止至8.1号
	if start.Unix() < DataStartTimestamp {
		start = time.Unix(DataStartTimestamp, 0)
	}
	if start.After(stop) {
		return "", "", errors.New("参数错误")
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
	//对应的es集群
	return strings.Join(esTableArr, ","), SureConnection(table), nil
}

//获取表命和对应的集群
func GetESTable(table string) (string, string) {
	//对应的es集群
	return table, SureConnection(table)
}

//获取表命和对应的集群
func GetESTableByDate(table, date string) (string, string) {
	//对应的es集群
	return fmt.Sprintf(table, date), SureConnection(table)
}
