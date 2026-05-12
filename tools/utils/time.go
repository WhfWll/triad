package utils

import (
	"sort"
	"strconv"
	"time"
)

const (
	DateTime  = "2006-01-02 15:04:05"
	DateOnly  = "2006-01-02"
	TimeOnly  = "15:04:05"
	YearMonth = "2006-01"
)

// DatetimeNumberStr 获取年月日时分秒数字化后组成的字符串
func DatetimeNumberStr() string {
	currentTime := time.Now()
	year := strconv.Itoa(currentTime.Year())
	month := currentTime.Format("01")
	day := strconv.Itoa(currentTime.Day())
	hour := strconv.Itoa(currentTime.Hour())
	minute := strconv.Itoa(currentTime.Minute())
	second := strconv.Itoa(currentTime.Second())
	return year + month + day + hour + minute + second
}

// WeekDateList 返回 最近一周的日列表、起始时间、数据库按日期统计查询时的格式化时间串
func WeekDateList() ([]string, string, string) {
	currentTime := time.Now()
	dateList := make([]string, 0)
	for i := 0; i < 7; i++ {
		dateList = append(dateList, currentTime.AddDate(0, 0, -i).Format(DateOnly))
	}
	sort.Strings(dateList)
	return dateList, currentTime.AddDate(0, 0, -6).Format(DateOnly), "%Y-%m-%d"
}

// MonthDateList 返回 最近一月的日列表、起始时间、数据库按日期统计查询时的格式化时间串
func MonthDateList() ([]string, string, string) {
	currentTime := time.Now()
	dateList := make([]string, 0)
	for i := 0; i < 31; i++ {
		dateList = append(dateList, currentTime.AddDate(0, 0, -i).Format(DateOnly))
	}
	sort.Strings(dateList)
	return dateList, currentTime.AddDate(0, 0, -30).Format(DateOnly), "%Y-%m-%d"
}

// YearDateList 返回 最近一年的月列表、起始时间、数据库按日期统计查询时的格式化时间串
func YearDateList() ([]string, string, string) {
	now := time.Now()
	currentTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	dateList := make([]string, 0)
	for i := 0; i < 12; i++ {
		dateList = append(dateList, currentTime.AddDate(0, -i, 0).Format(YearMonth))
	}
	sort.Strings(dateList)
	return dateList, currentTime.AddDate(0, -11, 0).Format(DateOnly), "%Y-%m"
}
