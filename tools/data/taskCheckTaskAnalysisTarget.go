package data

import (
	"regexp"
	"smart/tools/network"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TaskCheckTaskAnalysisTarget struct {
	TargetList           []string
	targetListMap        map[string]string
	ErrorTargetList      []string
	excludeTargetListMap map[string]string
	SkipUrlValidation    bool //是否跳过URL合法性校验
}

// AnalysisTarget 分析任务中的测试目标，解析出所需的目标列表
func (a *TaskCheckTaskAnalysisTarget) AnalysisTarget(checkTarget, excludeTarget string) {
	//根据任务中的测试目标，获取目标列表
	//taskTargetList := make([]string, 0)
	//a.targetListMap = make(map[string]string)
	//
	//checkTarget = strings.TrimSpace(checkTarget)
	//if strings.Contains(checkTarget, ",") {
	//	taskTargetList = strings.Split(checkTarget, ",")
	//} else {
	//	taskTargetList = strings.Split(checkTarget, "\n")
	//}

	taskTargetList := make([]string, 0)
	excludeTargetList := make([]string, 0)
	a.targetListMap = make(map[string]string)
	a.excludeTargetListMap = make(map[string]string)
	if strings.Contains(excludeTarget, ",") {
		excludeTargetList = strings.Split(excludeTarget, ",")
	} else {
		excludeTargetList = strings.Split(excludeTarget, "\n")
	}
	for _, one := range excludeTargetList {
		if one == "" {
			continue
		}
		a.excludeTargetListMap[strings.TrimSpace(one)] = "1"
	}
	checkTarget = strings.TrimSpace(checkTarget)
	if strings.Contains(checkTarget, ",") {
		taskTargetList = strings.Split(checkTarget, ",")
	} else {
		taskTargetList = strings.Split(checkTarget, "\n")
	}
	//解析目标中的ip段，分解为正常目标
	a.HandleIpSegment(taskTargetList)
}

// HandleIpSegment 处理ip段
func (a *TaskCheckTaskAnalysisTarget) HandleIpSegment(taskTargetList []string) {
	var invalidTargets []string
	for _, target := range taskTargetList {
		target = strings.TrimSpace(target)
		if a.FilterTarget(target) {
			continue
		}
		func(target string) {
			netmaskIpPattern := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+/\d+`)  //（粗略型）子网掩码型ip段正则
			crossbarIpPattern := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+-\d+`) //（粗略型）带横杠型ip段正则
			if netmaskIpPattern != nil && len(netmaskIpPattern.FindAllStringSubmatch(target, -1)) > 0 {
				if !network.IpSegmentTools.VerifyNetmaskIpSegment(target) {
					a.ErrorTargetList = append(a.ErrorTargetList, "子网掩码型ip段格式错误:"+target)
					return
				}
				targetSplit := strings.Split(target, "/")
				segmentNumber, _ := strconv.Atoi(targetSplit[len(targetSplit)-1])
				if segmentNumber < 16 {
					a.ErrorTargetList = append(a.ErrorTargetList, "子网掩码型网段不能小于21:"+target)
					return
				}
				ipList, err := network.IpSegmentTools.HandleNetmaskIpSegment(target)
				if err != nil {
					a.ErrorTargetList = append(a.ErrorTargetList, err.Error())
					return
				}
				for _, one := range ipList {
					a.AddTarget(one)
				}
			} else if crossbarIpPattern != nil && len(crossbarIpPattern.FindAllStringSubmatch(target, -1)) > 0 {
				if !network.IpSegmentTools.VerifyCrossbarIpSegment(target) {
					a.ErrorTargetList = append(a.ErrorTargetList, "带横杠型ip段格式错误:"+target)
					return
				}
				ipList, err := network.IpSegmentTools.HandleCrossbarIpSegment(target)
				if err != nil {
					a.ErrorTargetList = append(a.ErrorTargetList, err.Error())
					return
				}
				for _, one := range ipList {
					a.AddTarget(one)
				}
			} else {
				if _, err := network.IpSegmentTools.CheckSuccessUrl(target); err != nil {
					if a.SkipUrlValidation {
						a.AddTarget(target)
						return
					}
					// 收集非法目标
					invalidTargets = append(invalidTargets, target)
					return
				}
				a.AddTarget(target)
			}
		}(target)
	}

	// 如果有非法目标，统一添加到错误列表
	if len(invalidTargets) > 0 {
		msg := "存在非法测试目标：" + strings.Join(invalidTargets, "、")
		a.ErrorTargetList = append(a.ErrorTargetList, msg)
	}
}

// FilterTarget 判断目标是否需要过滤掉
func (a *TaskCheckTaskAnalysisTarget) FilterTarget(target string) bool {
	if target == "" {
		return true
	}
	return false
}

// AddTarget 添加目标到目标列表
func (a *TaskCheckTaskAnalysisTarget) AddTarget(target string) {
	if _, ok := a.targetListMap[target]; !ok {
		if _, ok := a.excludeTargetListMap[target]; ok {
			log.Println("目标被排除掉:" + target)
			return
		}

		a.targetListMap[target] = ""
		a.TargetList = append(a.TargetList, target)
	}
}
