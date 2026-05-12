package crons

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"smart/services"
	"smart/tools/enums"
	"strconv"
	"time"

	"github.com/mackerelio/go-osstat/network"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

type Monitor struct {
}

// MonitorHost 获取主机信息
func (m Monitor) MonitorHost() (*host.InfoStat, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return hostInfo, err
	}
	return hostInfo, nil
}

// MonitorCpu 获取CPU信息
func (m Monitor) MonitorCpu() ([]cpu.InfoStat, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return cpuInfo, err
	}
	return cpuInfo, nil
}

// MonitorCpuPercent 获取CPU使用率
func (m Monitor) MonitorCpuPercent() (int, error) {
	percent, _ := cpu.Percent(1*time.Second, false)
	var used int
	if len(percent) == 0 {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)
		used = int(stats.GCCPUFraction * 100)
	} else {
		used = int(math.Round(percent[0]))
	}
	return used, nil
}

// MonitorMemory 获取内存信息
func (m Monitor) MonitorMemory() (*mem.VirtualMemoryStat, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return memInfo, err
	}
	return memInfo, nil
}

func (m Monitor) MonitorMemoryPercent() (int, error) {
	memInfo, err := m.MonitorMemory()
	if err != nil {
		return 0, err
	}
	if memInfo.Total == 0 {
		return 0, nil
	}
	percent := int(math.Round(float64(memInfo.Used) / float64(memInfo.Total) * 100))
	return percent, nil
}

// MonitorDisk 获取磁盘分区信息
func (m Monitor) MonitorDisk() ([]disk.PartitionStat, error) {
	// 获取磁盘分区信息
	diskPartitions, err := disk.Partitions(true)
	if err != nil {
		fmt.Printf("Failed to get disk partitions: %v", err)
		return diskPartitions, err
	}
	fmt.Printf("Disk partitions: %+v\n", diskPartitions)
	for _, partition := range diskPartitions {
		// 获取每个磁盘分区的使用情况
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			fmt.Printf("Failed to get disk usage for %s: %v", partition.Mountpoint, err)
			continue
		}
		fmt.Printf("%s usage: %+v\n", partition.Mountpoint, usage)
	}
	return diskPartitions, nil
}

// MonitorRootDisk 获取根磁盘分区的使用情况
func (m Monitor) MonitorRootDisk() (*disk.UsageStat, error) {
	rootDisk, err := disk.Usage("/")
	if err != nil {
		return rootDisk, err
	}
	return rootDisk, nil
}

// CounterReceiveFlow 计算每秒的访问流量
func (m Monitor) CounterReceiveFlow() int {
	// 获取起始时间点的网络计数器
	beforeStats, err := network.Get()
	if err != nil {
		log.Errorf("get before network counter failed, err: %v\n", err)
		return 0
	}

	// 等待1秒钟
	time.Sleep(1 * time.Second)

	// 获取结束时间点的网络计数器
	afterStats, err := network.Get()
	if err != nil {
		log.Errorf("get after network counter failed, err: %v\n", err)
		return 0
	}

	// 计算起始时间点的接收、发送数据的总量
	beforeTxBytes := 0
	beforeRxBytes := 0
	for _, item := range beforeStats {
		if item.Name == "lo" {
			continue
		}
		beforeTxBytes += int(item.TxBytes)
		beforeRxBytes += int(item.RxBytes)
	}
	// 计算结束时间点的接收、发送数据的总量
	afterTxBytes := 0
	afterRxBytes := 0
	for _, item := range afterStats {
		if item.Name == "lo" {
			continue
		}
		afterTxBytes += int(item.TxBytes)
		afterRxBytes += int(item.RxBytes)
	}

	// 计算接收流量
	receiveFlow := afterRxBytes - beforeRxBytes
	//fmt.Println(receiveFlow)
	// 计算发送流量
	//sendFlow := afterTxBytes - beforeTxBytes
	//fmt.Println(sendFlow)
	return receiveFlow
}

// GetNodeTime 获取时间轴中的节点时间
func (m Monitor) GetNodeTime() string {
	currentTime := time.Now()
	hour := strconv.Itoa(currentTime.Hour())
	minute := strconv.Itoa(currentTime.Minute())
	if len(hour) == 1 {
		hour = "0" + hour
	}
	if len(minute) == 1 {
		minute = "0" + minute
	}
	xData := hour + ":" + minute + ":" + "00"
	return xData
}

func SystemMonitor() {
	log.Info("系统监控中...")
	var (
		monitor       Monitor
		mapSetService services.MapSet
		ctx           = context.Background()
	)
	nodeTime := monitor.GetNodeTime()

	//cpu
	go func() {
		var cpuMapSet services.MonitorCpuMemoryMapSet
		cpuPercent, err := monitor.MonitorCpuPercent()
		if err != nil {
			log.Error("SystemMonitor MonitorCpuPercent err：" + err.Error())
			return
		}
		HandleMonitorWarn(ctx, enums.MonitorCpu, cpuPercent)
		cpuItem := services.MonitorCpuMemoryItem{
			XData: nodeTime,
			YData: cpuPercent,
		}
		//获取cpu的ObjValue数据
		objValueStr, err := mapSetService.GetMapValue(ctx, enums.MonitorCpuMapSetObjKey)
		if err != nil {
			log.Error("SystemMonitor get cpu objValue err：" + err.Error())
			return
		}
		if objValueStr == "" {
			cpuMapSet.List = append(cpuMapSet.List, cpuItem)
		} else {
			if err = json.Unmarshal([]byte(objValueStr), &cpuMapSet); err != nil {
				log.Error("SystemMonitor parse cpu objValue err：" + err.Error())
				return
			}
			if len(cpuMapSet.List) < 30 {
				cpuMapSet.List = append(cpuMapSet.List, cpuItem)
			} else {
				cpuMapSet.List = cpuMapSet.List[len(cpuMapSet.List)-29 : len(cpuMapSet.List)]
				cpuMapSet.List = append(cpuMapSet.List, cpuItem)
			}
		}
		//更新后的cpu数据存入mapSet
		cpuByte, err := json.Marshal(cpuMapSet)
		if err != nil {
			log.Error("SystemMonitor Marshal cpu objValue err：" + err.Error())
			return
		}
		if err = mapSetService.Create(ctx, enums.MonitorCpuMapSetObjKey, string(cpuByte), enums.MonitorCpuMapSetContent); err != nil {
			log.Error("SystemMonitor save cpu objValue err：" + err.Error())
			return
		}
	}()

	//内存
	go func() {
		var memoryMapSet services.MonitorCpuMemoryMapSet
		memPercent, err := monitor.MonitorMemoryPercent()
		if err != nil {
			log.Error("SystemMonitor MonitorMemoryPercent err：" + err.Error())
			return
		}
		HandleMonitorWarn(ctx, enums.MonitorMemory, memPercent)
		memoryItem := services.MonitorCpuMemoryItem{
			XData: nodeTime,
			YData: memPercent,
		}
		//获取memory的ObjValue数据
		memoryStr, err := mapSetService.GetMapValue(ctx, enums.MonitorMemoryMapSetObjKey)
		if err != nil {
			log.Error("SystemMonitor get memory objValue err：" + err.Error())
			return
		}
		if memoryStr == "" {
			memoryMapSet.List = append(memoryMapSet.List, memoryItem)
		} else {
			if err = json.Unmarshal([]byte(memoryStr), &memoryMapSet); err != nil {
				log.Error("SystemMonitor parse memory objValue err：" + err.Error())
				return
			}
			if len(memoryMapSet.List) < 30 {
				memoryMapSet.List = append(memoryMapSet.List, memoryItem)
			} else {
				memoryMapSet.List = memoryMapSet.List[len(memoryMapSet.List)-29 : len(memoryMapSet.List)]
				memoryMapSet.List = append(memoryMapSet.List, memoryItem)
			}
		}
		//更新后的memory数据存入mapSet
		memoryByte, err := json.Marshal(memoryMapSet)
		if err != nil {
			log.Error("SystemMonitor Marshal memory objValue err：" + err.Error())
			return
		}
		if err = mapSetService.Create(ctx, enums.MonitorMemoryMapSetObjKey, string(memoryByte), enums.MonitorMemoryMapSetContent); err != nil {
			log.Error("SystemMonitor save memory objValue err：" + err.Error())
			return
		}
	}()

	//硬盘
	go func() {
		var diskMapSet services.MonitorDiskMapSet
		diskInfo, err := monitor.MonitorRootDisk()
		if err != nil {
			log.Error("SystemMonitor MonitorRootDisk err：" + err.Error())
			return
		}
		diskFree := diskInfo.Free
		diskUsed := diskInfo.Used
		diskTotal := diskFree + diskUsed
		diskFreePercent := int(math.Round(float64(diskFree) / float64(diskTotal) * 100))
		diskUsedPercent := 100 - diskFreePercent
		diskMapSet.Free = int(math.Round(float64(diskFree) / 1024 / 1024 / 1024))
		diskMapSet.Used = int(math.Round(float64(diskUsed) / 1024 / 1024 / 1024))
		diskMapSet.Total = diskMapSet.Used + diskMapSet.Free
		diskMapSet.FreePercent = diskFreePercent
		diskMapSet.UsedPercent = diskUsedPercent
		HandleMonitorWarn(ctx, enums.MonitorDisk, diskUsedPercent)
		diskByte, err := json.Marshal(diskMapSet)
		if err != nil {
			log.Error("SystemMonitor json Marshal err：" + err.Error())
			return
		}
		if err = mapSetService.Create(ctx, enums.MonitorDiskMapSetObjKey, string(diskByte), enums.MonitorDiskMapSetContent); err != nil {
			log.Error("SystemMonitor save disk objValue err：" + err.Error())
			return
		}
	}()

	//流量
	go func() {
		flow := monitor.CounterReceiveFlow()
		HandleMonitorWarn(ctx, enums.MonitorFlow, flow)
	}()
}

// HandleMonitorWarn 处理系统监控告警的信息
func HandleMonitorWarn(ctx context.Context, warnModel string, value int) {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.MonitorWarnMapSetObjKey)
	if err != nil {
		log.Error("SystemConfig MonitorWarn GetMapValue err：" + err.Error())
		return
	}
	if objValueStr == "" {
		log.Error("SystemConfig MonitorWarn objValue is null")
		return
	}
	var monitorWarnMapSet services.MonitorWarnMapSet
	if err = json.Unmarshal([]byte(objValueStr), &monitorWarnMapSet); err != nil {
		log.Error("SystemConfig MonitorWarn Unmarshal err：" + err.Error())
		return
	}
	if monitorWarnMapSet.IsOpen == enums.ConfigClose {
		return
	}
	var (
		content       string
		contentSuffix = "超过"
	)
	switch warnModel {
	case enums.MonitorCpu:
		if value > monitorWarnMapSet.CpuWarn {
			content = enums.MonitorCpu + contentSuffix + strconv.Itoa(monitorWarnMapSet.CpuWarn) + "%"
		}
	case enums.MonitorMemory:
		if value > monitorWarnMapSet.MemoryWarn {
			content = enums.MonitorMemory + contentSuffix + strconv.Itoa(monitorWarnMapSet.MemoryWarn) + "%"
		}
	case enums.MonitorDisk:
		if value > monitorWarnMapSet.DiskWarn {
			content = enums.MonitorDisk + contentSuffix + strconv.Itoa(monitorWarnMapSet.DiskWarn) + "%"
		}
	case enums.MonitorFlow:
		if value > monitorWarnMapSet.FlowWarn*1024*1024 {
			content = enums.MonitorFlow + contentSuffix + strconv.Itoa(monitorWarnMapSet.FlowWarn) + "M"
		}
	}
	if content == "" {
		return
	}
	var messageService services.SystemMessage
	if err := messageService.SystemMessageAdd(ctx, content, enums.MessageTypeWarn, 0, enums.MessageStatusUnread); err != nil {
		log.Error("SystemMessageAdd err：" + err.Error())
		return
	}
	var logAuditService services.LogAudit
	if err := logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeWarn, content, "admin", "127.0.0.1"); err != nil {
		log.Error("LogAuditAdd err：" + err.Error())
		return
	}
}
