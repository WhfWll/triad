package data

import (
	"errors"
	"fmt"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"
)

// 渗透任务 下次运行时间计算
type TaskExecNextTime struct {
	StartTime          time.Time `form:"startTime" json:"startTime"`                   // 开始时间
	CyclePlanningType  int       `form:"cyclePlanningType" json:"cyclePlanningType"`   // 周期计划类型 1周 2月
	CyclePlanningValue int       `form:"cyclePlanningValue" json:"cyclePlanningValue"` // 周期计划类型 具体时间 （其中周：1周一 2周二... 7周日） （其中月：1一号 2二号...）
	CyclePlanningHour  string    `form:"cyclePlanningHour" json:"cyclePlanningHour"`   // 周期计划类型 具体小时与分钟 00:01 ｜ 00:20等
	EndTime            time.Time `form:"endTime" json:"endTime"`                       // 结束时间
}

// 测试
func TestNextTime() {
	// 即时执行
	testCall(1, "", 0, 0, "", "")

	// 定时执行
	testCall(2, "2023-10-01 12:12:12", 0, 0, "", "")

	// 周期执行
	// 每周执行
	fmt.Println("每周执行")
	// 15   16  17  18  19  20  21  22  23  24
	// 4    5   6   7   1   2   3   4   5   6
	// 下周1
	testCall(3, "", 1, 1, "12:12", "2023-10-01 12:12:12")
	// 下周2
	testCall(3, "", 1, 2, "12:12", "2023-10-01 12:12:12")
	// 下周3
	testCall(3, "", 1, 3, "12:12", "2023-10-01 12:12:12")
	// 下周4
	testCall(3, "", 1, 4, "12:12", "2023-10-01 12:12:12")
	// 下周5
	testCall(3, "", 1, 5, "12:12", "2023-10-01 12:12:12")
	// 下周6
	testCall(3, "", 1, 6, "12:12", "2023-10-01 12:12:12")
	// 下周7
	testCall(3, "", 1, 7, "12:12", "2023-10-01 12:12:12")

	// 每月执行
	fmt.Println("每月执行")
	// 1号
	testCall(3, "", 2, 1, "12:12", "2023-10-01 12:12:12")
	// 10号
	testCall(3, "", 2, 10, "12:12", "2023-10-01 12:12:12")
	// 15号
	testCall(3, "", 2, 15, "12:12", "2023-10-01 12:12:12")
	// 20号
	testCall(3, "", 2, 20, "12:12", "2023-10-01 12:12:12")
	// 30号
	testCall(3, "", 2, 30, "12:12", "2023-10-01 12:12:12")
	// 31号
	testCall(3, "", 2, 31, "12:12", "2023-10-01 12:12:12")
}
func testCall(execType int, startTime string, planType, planValue int, hour string, endTime string) {
	var taskExecNextTime TaskExecNextTime
	taskExecNextTime.StartTime, _ = time.ParseInLocation(utils.DateTime, startTime, time.Local)
	taskExecNextTime.CyclePlanningType = planType
	taskExecNextTime.CyclePlanningValue = planValue
	taskExecNextTime.CyclePlanningHour = hour
	taskExecNextTime.EndTime, _ = time.ParseInLocation(utils.DateTime, endTime, time.Local)
	nextRuntime, err := taskExecNextTime.Compute(execType)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(nextRuntime)
}

func (execNextTime TaskExecNextTime) Compute(taskType int) (time.Time, error) {

	switch taskType {
	case enums.TaskExecTypeImmediate: // 即时执行
		// 即时执行任务，下次运行时间为当前时间，说明任务需要立即执行
		return time.Now(), nil
	case enums.TaskExecTypeTiming: // 定时执行
		// 定时执行，下次运行时间为指定的时间
		if execNextTime.StartTime.IsZero() {
			return time.Time{}, errors.New("渗透任务 - 下次运行时间计算 - 定时执行任务开始时间不可为空")
		}
		return execNextTime.StartTime, nil
	case enums.TaskExecTypeCycle: // 周期执行
		// 周期执行，计算下一次周期执行的时间
		if execNextTime.CyclePlanningType == 0 {
			return time.Time{}, errors.New("渗透任务 - 下次运行时间计算 - 周期执行任务执行类型未知")
		}
		if execNextTime.CyclePlanningValue == 0 {
			return time.Time{}, errors.New("渗透任务 - 下次运行时间计算 - 周期执行任务执行数值未知")
		}
		if execNextTime.CyclePlanningHour == "" {
			return time.Time{}, errors.New("渗透任务 - 下次运行时间计算 - 周期执行任务执行具体时间未知")
		}
		if execNextTime.EndTime.IsZero() {
			return time.Time{}, errors.New("渗透任务 - 下次运行时间计算 - 周期执行任务结束时间不可为空")
		}

		// Sunday、Monday、Tuesday、Wednesday、Thursday、Friday和Saturday
		// 星期日、星期一、星期二、星期三、星期四、星期五和星期六

		// 是否已过结束时间
		if time.Now().Unix() >= execNextTime.EndTime.Unix() {
			return time.Time{}, errors.New("渗透任务 - 下次运行时间计算 - 已超过结束时间")
		}

		// 每周执行一次
		if execNextTime.CyclePlanningType == enums.TaskExecTypeCycleTypeWeek {
			currentWeekEn := time.Now().Weekday()
			var currentWeekNum int // 计算当前周几
			switch currentWeekEn.String() {
			case "Monday": // 星期一
				currentWeekNum = 1
			case "Tuesday": // 星期二
				currentWeekNum = 2
			case "Wednesday": // 星期三
				currentWeekNum = 3
			case "Thursday": // 星期四
				currentWeekNum = 4
			case "Friday": // 星期五
				currentWeekNum = 5
			case "Saturday": // 星期六
				currentWeekNum = 6
			case "Sunday": // 星期日
				currentWeekNum = 7
			}

			// 当前周大于计划周
			if currentWeekNum > execNextTime.CyclePlanningValue {
				// 计算：距离下次计划周几还差几天
				return execNextTime.nextWeekTime(currentWeekNum), nil
			}

			// 当前周等于计划周
			if currentWeekNum == execNextTime.CyclePlanningValue {
				// 计算小时是否大于
				hourMin := strings.Split(execNextTime.CyclePlanningHour, ":")
				hourInt, _ := strconv.Atoi(hourMin[0])
				minInt, _ := strconv.Atoi(hourMin[1])
				if time.Now().Hour() > hourInt {
					// 计算下周几
					return execNextTime.nextWeekTime(currentWeekNum), nil
				} else if time.Now().Hour() == hourInt {
					// 计算分钟
					if time.Now().Minute() > minInt || time.Now().Minute() == minInt {
						// 计算下周几
						return execNextTime.nextWeekTime(currentWeekNum), nil
					} else {
						// 计算当天时间
						return time.ParseInLocation(utils.DateTime, time.Now().Format(utils.DateOnly)+" "+execNextTime.CyclePlanningHour+":00", time.Local)
					}
				} else {
					// 当前时间小于计划时间，计算当前时间
					return time.ParseInLocation(utils.DateTime, time.Now().Format(utils.DateOnly)+" "+execNextTime.CyclePlanningHour+":00", time.Local)
				}
			}

			// 当前周小于计划周
			if currentWeekNum < execNextTime.CyclePlanningValue {
				diffDay := execNextTime.CyclePlanningValue - currentWeekNum
				nextDateTime := time.Now().AddDate(0, 0, diffDay).Format(utils.DateOnly) + " " + execNextTime.CyclePlanningHour + ":00"
				return time.ParseInLocation(utils.DateTime, nextDateTime, time.Local)
			}
		}

		// 每月执行一次
		if execNextTime.CyclePlanningType == enums.TaskExecTypeCycleTypeMonth {
			currentDay := time.Now().Day()

			// 当前号 大于 计划号
			if currentDay > execNextTime.CyclePlanningValue {
				// 计算下月目标号
				return execNextTime.nextMonthDayTime()
			}

			// 当前号 等于 计划号
			if currentDay == execNextTime.CyclePlanningValue {
				// 计算时间是否大于
				hourMin := strings.Split(execNextTime.CyclePlanningHour, ":")
				hourInt, _ := strconv.Atoi(hourMin[0])
				minInt, _ := strconv.Atoi(hourMin[1])
				if time.Now().Hour() > hourInt {
					return execNextTime.nextMonthDayTime()
				} else if time.Now().Hour() == hourInt {
					// 计算分钟
					if time.Now().Minute() > minInt || time.Now().Minute() == minInt {
						// 计算下周几
						return execNextTime.nextMonthDayTime()
					} else {
						// 计算当天时间
						return time.ParseInLocation(utils.DateTime, time.Now().Format(utils.DateOnly)+" "+execNextTime.CyclePlanningHour+":00", time.Local)
					}
				} else {
					// 当前时间小于计划时间，计算当前时间
					return time.ParseInLocation(utils.DateTime, time.Now().Format(utils.DateOnly)+" "+execNextTime.CyclePlanningHour+":00", time.Local)
				}
			}

			// 当前号 小于 计划号
			if currentDay < execNextTime.CyclePlanningValue {
				planValue := strconv.Itoa(execNextTime.CyclePlanningValue)
				if execNextTime.CyclePlanningValue < 10 {
					planValue += "0" + strconv.Itoa(execNextTime.CyclePlanningValue)
				}

				currentDate := time.Now().Format(utils.DateOnly)
				currentDateSlice := strings.Split(currentDate, "-")
				return time.ParseInLocation(utils.DateTime, currentDateSlice[0]+"-"+currentDateSlice[1]+"-"+planValue+" "+execNextTime.CyclePlanningHour+":00", time.Local)
			}
		}
	}
	return time.Time{}, nil
}

func (execNextTime TaskExecNextTime) nextWeekTime(currentWeekNum int) time.Time {
	nextDay := 7 - currentWeekNum + execNextTime.CyclePlanningValue
	nextDate := time.Now().AddDate(0, 0, nextDay).Format(utils.DateOnly)
	nextDateTime := nextDate + " " + execNextTime.CyclePlanningHour + ":00"
	finalTime, _ := time.ParseInLocation(utils.DateTime, nextDateTime, time.Local)
	return finalTime
}

func (execNextTime TaskExecNextTime) nextMonthDayTime() (time.Time, error) {
	target := execNextTime.CyclePlanningValue

	var finalDateTime time.Time

	// 下月
	nextMonth := time.Now().AddDate(0, 1, 0)
	if nextMonth.Day() > target {
		finalDateTime = nextMonth.AddDate(0, 0, target-nextMonth.Day())
	} else if nextMonth.Day() == target {
		finalDateTime = nextMonth
	} else {
		// 保存当前月份，这里有坑，如果是2月27日 但是target是31 那么这里的月份就会到3月份，记录下，然后验证
		month := nextMonth.Month()
		nextMonth = nextMonth.AddDate(0, 0, target-nextMonth.Day())
		if month != nextMonth.Month() {
			// 说明月份已发生变化，时间有问题，不能处理
			return time.Time{}, errors.New("下个月无此执行日期")
		}
		finalDateTime = nextMonth
	}

	return time.ParseInLocation(utils.DateTime, finalDateTime.Format(utils.DateOnly)+" "+execNextTime.CyclePlanningHour+":00", time.Local)
}
