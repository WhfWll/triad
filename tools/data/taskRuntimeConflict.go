package data

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"
)

type TaskRuntimeCheck struct {
	StartTime          string   `form:"startTime" json:"startTime"`                   // 开始时间
	CyclePlanningType  int      `form:"cyclePlanningType" json:"cyclePlanningType"`   // 周期计划类型 1周 2月
	CyclePlanningValue int      `form:"cyclePlanningValue" json:"cyclePlanningValue"` // 周期计划类型 具体周期 （其中周：1周一 2周二... 7周日） （其中月：1一号 2二号...）
	CyclePlanningHour  string   `form:"cyclePlanningHour" json:"cyclePlanningHour"`   // 周期计划类型 具体小时与分钟 00:01 ｜ 00:20等
	EndTime            string   `form:"endTime" json:"endTime"`                       // 结束时间
	RuntimePeriod      []string `form:"runtimePeriod" json:"runtimePeriod" binding:"required"`
}

// 校验运行时间数据是否非法
func (check *TaskRuntimeCheck) CheckRuntime(execType int) (err error) {
	// 运行时间段 校验
	if len(check.RuntimePeriod) == 0 {
		return errors.New("运行时间段不可为空")
	}

	// 校验运行时间段是否与运行方式有冲突
	// 可能冲突1 定时执行的执行时间 没有在可运行时间段内
	// 可能冲突2 周期任务开始时间 没有在可运行的时间段内
	// 校验执行类型参数
	switch execType {
	case 1: // 即时执行
	case 2: // 定时执行
		if check.StartTime == "" {
			return errors.New("定时执行任务，开始时间不可为空")
		}
		// 校验计划时间是否早于当前时间
		startTime, err := time.ParseInLocation(utils.DateTime, check.StartTime, time.Local)
		if err != nil {
			return errors.New("执行方式时间格式无法识别")
		}
		if startTime.Before(time.Now()) {
			return errors.New("计划时间不能早于当前时间")
		}
		// 可能冲突1 定时执行的执行时间 没有在可运行时间段内
		return check.checkConflict(check.StartTime, check.RuntimePeriod)
	case 3: // 周期执行
		// 执行计划类型
		if check.CyclePlanningType != enums.TaskExecTypeCycleTypeWeek && check.CyclePlanningType != enums.TaskExecTypeCycleTypeMonth {
			return errors.New("周期执行任务，周期计划类型不可为空")
		}

		// 执行计划类型对应的值
		if check.CyclePlanningValue == 0 {
			return errors.New("周期执行任务，周期计划值不可为空")
		}

		// 执行计划对应时间 格式为 00:00
		if check.CyclePlanningHour == "" {
			return errors.New("周期执行任务，周期计划值不可为空")
		}
		CyclePlanningHourSlice := strings.Split(check.CyclePlanningHour, ":")
		if len(CyclePlanningHourSlice) != 2 {
			return errors.New("周期执行任务，周期计划值非法")
		}
		_, err = strconv.Atoi(CyclePlanningHourSlice[0])
		if err != nil {
			return errors.New("周期执行任务，周期计划值非法")
		}
		_, err = strconv.Atoi(CyclePlanningHourSlice[1])
		if err != nil {
			return errors.New("周期执行任务，周期计划值非法")
		}
		// 可能冲突2 周期任务开始时间 没有在可运行的时间段内
		// 由于周期任务的开始时间为CyclePlanningHour 且格式为 00:00 所以校验冲突时，可将其组织为正常的一个时间
		err = check.checkConflict("2023-06-15 "+check.CyclePlanningHour+":00", check.RuntimePeriod)
		if err != nil {
			return err
		}
		// 周期计划 - 结束时间
		if check.EndTime == "" {
			return errors.New("周期执行任务，结束时间不可为空")
		}
		_, err = time.Parse(utils.DateTime, check.EndTime)
		if err != nil {
			return errors.New("周期执行任务，结束时间格式无法识别")
		}
	default:
		return errors.New("未知的执行方式")
	}

	return nil
}

// 校验运行时间与运行时间段冲突
func (check *TaskRuntimeCheck) checkConflict(timeStr string, runtimeList []string) error {
	if timeStr == "" {
		return errors.New("校验运行时间段是否冲突，timeStr不可为空")
	}
	if len(runtimeList) == 0 {
		return errors.New("校验运行时间段是否冲突，runtimeList不可为空")
	}
	// 计算开始时间 开始时间的当天0点至开始时间的秒数 如2023-06-14 12:00:00 那么仅计算 0点到12:00:00的秒数
	startTime, err := time.Parse(utils.DateTime, timeStr)
	if err != nil {
		return errors.New("执行方式时间格式无法识别")
	}
	startUnix := startTime.Unix()
	startDateTime, _ := time.Parse(utils.DateOnly, startTime.Format(utils.DateOnly))
	startDateUnix := startDateTime.Unix()
	second := int(startUnix - startDateUnix)

	flag := false
	for _, item := range runtimeList {
		runtimeStartSecond, runtimeEndSecond, err := check.getRuntimeSecond(item)
		if err != nil {
			return err
		}
		if second >= runtimeStartSecond && second < runtimeEndSecond {
			flag = true
		}
	}

	// 不满足
	if flag == false {
		return errors.New("运行时间段与执行方式时间冲突，请确认")
	}

	return nil
}

func (check *TaskRuntimeCheck) getRuntimeSecond(runtimePeriod string) (int, int, error) {
	itemSlice := strings.Split(runtimePeriod, "-")
	if len(itemSlice) != 2 {
		return 0, 0, errors.New("运行时间段格式存在非法")
	}

	// 开始
	firstSlice := strings.Split(itemSlice[0], ":")
	if len(firstSlice) != 2 {
		return 0, 0, errors.New("运行时间段格式存在非法")
	}
	// 开始小时
	firstHour, err := strconv.Atoi(firstSlice[0])
	if err != nil {
		return 0, 0, errors.New("运行时间段格式存在非法")
	}
	// 开始分
	firstMin, err := strconv.Atoi(firstSlice[1])
	if err != nil {
		return 0, 0, errors.New("运行时间段格式存在非法")
	}
	// 开始时间转换秒
	startSecond := firstHour*3600 + firstMin*60

	// 结束
	endSlice := strings.Split(itemSlice[1], ":")
	if len(endSlice) != 2 {
		return 0, 0, errors.New("运行时间段格式存在非法")
	}
	// 结束小时
	endHour, err := strconv.Atoi(endSlice[0])
	if err != nil {
		return 0, 0, errors.New("运行时间段格式存在非法")
	}
	// 注意：这里有可能是23:00-00:00 所有需要特殊处理
	if firstHour > 0 && endHour == 0 {
		endHour = firstHour + 1
	}
	// 结束分
	endMin, err := strconv.Atoi(endSlice[1])
	if err != nil {
		return 0, 0, errors.New("运行时间段格式存在非法")
	}
	// 结束时间转换秒
	endSecond := endHour*3600 + endMin*60

	return startSecond, endSecond, nil
}

// 校验运行时间段是否处于当前时间，是否允许运行
func (check *TaskRuntimeCheck) CheckIsAllowRunning(runtimePeriod []string) bool {
	// 默认不允许运行
	isAllowRun := false
	currentTotalMin := (time.Now().Hour() * 60) + time.Now().Minute()

	for _, runtime := range runtimePeriod {
		// 校验 00:00-01:00 一组格式是否正确
		startAndEnd := strings.Split(runtime, "-")
		if len(startAndEnd) != 2 {
			log.Error("运行时间段非法：" + runtime)
			continue
		}

		/********** 开始时间校验 ************/
		// 当天0点到允许时间的分钟数
		allowStartTotalMin, err := check.jisuanTotalMin(startAndEnd[0])
		if err != nil {
			log.Error("运行时间段非法：" + runtime)
			continue
		}

		/********** 结束时间校验 ************/
		// 当天0点到允许时间的分钟数
		allowEndTotalMin, err := check.jisuanTotalMin(startAndEnd[1])
		if err != nil {
			log.Error("运行时间段非法：" + runtime)
			continue
		}
		// 如果是23:00-00:00 则需要给结束时间+60分钟
		if allowEndTotalMin < allowStartTotalMin {
			allowEndTotalMin = allowStartTotalMin + 60
		}

		/********** 验证当前时间是否在允许的时间段内 ************/
		//fmt.Println("当前时间", currentTotalMin)
		//fmt.Println("开始时间", allowStartTotalMin)
		//fmt.Println("结束时间", allowEndTotalMin)
		//fmt.Println("---------------------------")
		if currentTotalMin >= allowStartTotalMin && currentTotalMin <= allowEndTotalMin {
			isAllowRun = true
		}
	}

	return isAllowRun
}

// @param 00:00 小时:分
// @return 0点到param一共多少分钟
func (check *TaskRuntimeCheck) jisuanTotalMin(hourMin string) (int, error) {
	// 校验时间 00:00 是否正确
	startHourAndMin := strings.Split(hourMin, ":")
	if len(startHourAndMin) != 2 {
		return 0, errors.New("参数非法")
	}
	// 校验时间是否可以转换未int 小时
	allowStartHour, err := strconv.Atoi(startHourAndMin[0])
	if err != nil {
		return 0, errors.New("小时转换失败 err = " + err.Error())
	}
	// 校验开始时间是否可以转换未int 分
	allowStartMin, err := strconv.Atoi(startHourAndMin[1])
	if err != nil {
		return 0, errors.New("分钟转换失败 err = " + err.Error())
	}
	// 当天0点到允许时间的分钟数
	allowStartTotalMin := (allowStartHour * 60) + allowStartMin

	return allowStartTotalMin, nil
}
