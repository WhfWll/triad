package time

import (
	"fmt"
	"strings"
	"time"
)

// 获取当天0点的日期时间
func GetTodyZeroTime() string {
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Format("2006-01-02 15:04:05")
}

// 获取当天0点的日期时间（格式化）
func GetTodyZeroTimeFormat(timeFormat string) string {
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Format(timeFormat)
}

// 获取当月1号0点的日期时间（格式化）
func GetMonthZeroTimeFormat(timeFormat string) string {
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, currentTime.Location()).Format(timeFormat)
}

// 获取当前日期时间
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取n天前的时间
func GetBeforeDayTime(n int) string {
	return time.Now().AddDate(0, 0, -n).Format("2006-01-02 15:04:05")
}

// 获取n个小时前的时间
func GetBeforeHourTime(n int, layout string) string {
	return time.Now().Add(time.Hour * time.Duration(-n)).Format(layout)
}

func GetBeforeDayTimeFormat(n int, timeFormat string) string {
	return time.Now().AddDate(0, 0, -n).Format(timeFormat)
}

// 获取n分钟之前的时间
func GetBeforeMinuteTime(n int) string {
	return time.Now().Add(time.Minute * time.Duration(-n)).Format("2006-01-02 15:04:05")
}

// 获取n秒之前的时间
func GetBeforeSecondTime(n int) string {
	return time.Now().Add(time.Second * time.Duration(-n)).Format("2006-01-02 15:04:05")
}

// 字符串转时间
func TimeStringToTime(timeString, layout string) (time.Time, error) {
	local, err := time.LoadLocation("Asia/Shanghai") //时间
	if err != nil {
		return time.Time{}, err
	}
	sTime, err := time.ParseInLocation(layout, timeString, local)
	return sTime, err
}

// TimeStringToInt64 字符串转时间戳
func TimeStringToInt64(timeString, layout string) int64 {
	local, _ := time.LoadLocation("Asia/Shanghai") //时间
	sTime, _ := time.ParseInLocation(layout, timeString, local)
	return sTime.Unix()
}

// CalculationESTime 计算web异常预警检索条件 开始时间
func CalculationESTime(currentTime time.Time, reduce, add string) int64 {
	timeT := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, currentTime.Location())
	endReduce2, _ := time.ParseDuration(reduce)
	calTime := timeT.Add(endReduce2)
	if add != "" {
		endAdd59, _ := time.ParseDuration(add)
		calTime = calTime.Add(endAdd59)
	}
	return calTime.UnixMilli()
}

// UTCToBeijingTime UTC时间转北京时间
func UTCToBeijingTime(utcTimeStr string) (string, error) {
	// 解析UTC时间
	utcTime, err := time.Parse(time.RFC3339, utcTimeStr)
	if err != nil {
		return "", err
	}

	if utcTime.IsZero() {
		return "", nil
	}

	// 转换为北京时间 (UTC+8)
	beijingTime := utcTime.Add(8 * time.Hour)

	// 格式化输出
	return beijingTime.Format("2006-01-02 15:04:05"), nil
}

// 智能转换时间，自动处理不同精度
func SmartConvertToBeijing(timeStr string) (string, error) {
	if timeStr == "" {
		return "", nil
	}

	// 尝试不同的时间格式
	var parsedTime time.Time
	var err error

	// 检查是否包含小数秒
	if strings.Contains(timeStr, ".") {
		// 尝试纳秒精度
		parsedTime, err = time.Parse(time.RFC3339Nano, timeStr)
		if err != nil {
			return "", err
		}
	} else {
		// 尝试标准 RFC3339 格式
		parsedTime, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return "", err
		}
	}

	// 零值检查
	if parsedTime.IsZero() || (parsedTime.Year() <= 1 && parsedTime.Month() <= 1 && parsedTime.Day() <= 1) {
		return "", nil
	}

	// 转换为北京时间
	beijingTime := parsedTime.Add(8 * time.Hour)

	return beijingTime.Format("2006-01-02 15:04:05"), nil
}

// 生成指定时间的 RFC3339 格式（带时区偏移）
func GenerateRFC3339WithOffset(year, month, day, hour, min, sec int) string {
	// 创建时间对象（UTC时区）
	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)

	// 格式化为带时区偏移的格式
	return t.Format("2006-01-02T15:04:05-07:00")
}

// 生成当前时间的 RFC3339 格式（带时区偏移）
func GenerateCurrentRFC3339WithOffset() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05-07:00")
}

// FormatDuration 将一个以秒为单位的时长转换为易读的中文格式字符串。
func FormatDuration(seconds int64) string {
	if seconds <= 0 {
		return "0秒"
	}
	// 将秒数转换为 time.Duration 类型以便进行运算
	duration := time.Duration(seconds) * time.Second
	const (
		secondsPerMinute = 60
		secondsPerHour   = 60 * secondsPerMinute
		secondsPerDay    = 24 * secondsPerHour
		secondsPerMonth  = 30 * secondsPerDay  // 近似值
		secondsPerYear   = 365 * secondsPerDay // 近似值
	)

	// 从年开始计算
	years := seconds / secondsPerYear
	seconds %= secondsPerYear

	months := seconds / secondsPerMonth
	seconds %= secondsPerMonth

	days := seconds / secondsPerDay
	seconds %= secondsPerDay

	hours := seconds / secondsPerHour
	seconds %= secondsPerHour

	minutes := seconds / secondsPerMinute
	seconds %= secondsPerMinute

	parts := []string{}

	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d年", years))
	}
	if months > 0 {
		parts = append(parts, fmt.Sprintf("%d个月", months))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d天", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d小时", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d分钟", minutes))
	}
	// 只有在总时长大于 0，并且不是由更大的单位完全组成时，才显示秒
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d秒", seconds))
	}
	if len(parts) == 0 && duration > 0 {
		// 这一步理论上不会发生，但作为安全检查
		return "不足1秒"
	}
	return strings.Join(parts, "")
}
