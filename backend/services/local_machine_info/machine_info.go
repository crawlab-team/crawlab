package local_machine_info

import (
	"crawlab/utils"
	"github.com/apex/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type MachineInfo struct {
	CpuUsed   float64 // cpu 整体使用率
	CpuNums   int     // 物理cpu数量
	MemTotal  float64 // 内存总量  unit: Mb
	MemFree   float64 // 内存空闲  unit: Mb
	MemUsed   float64 // 内存使用率
	DhUsed    float64 // 磁盘使用率
	DhFree    float64 // 磁盘空闲  unit: G
	DhTotal   float64 // 磁盘总量  unit: G
	LoadAvg1  float64 // 机器cpu负载1分钟统计
	LoadAvg5  float64 // 机器cpu负载5分钟统计
	LoadAvg15 float64 // 机器cpu负载15分钟统计
}

func CollectMachineInfo() MachineInfo {
	curMem, _ := mem.VirtualMemory()
	curCpuNums, _ := cpu.Counts(false)
	curCpuPercent, _ := cpu.Percent(time.Second, false)
	curDh, _ := disk.Usage("/")
	loadInfo, _ := load.Avg()

	if len(curCpuPercent) == 0 {
		log.Warn("unable get cpu usage")
		curCpuPercent = append(curCpuPercent, -1.0)
	}

	machineUsageInfo := MachineInfo{
		CpuUsed:   utils.Decimal(curCpuPercent[0]),
		CpuNums:   curCpuNums,
		MemTotal:  float64(curMem.Total / 1024 / 1024),
		MemFree:   float64(curMem.Available / 1024 / 1024),
		MemUsed:   utils.Decimal(curMem.UsedPercent),
		DhTotal:   float64(curDh.Total / 1024 / 1024 / 1024),
		DhFree:    float64(curDh.Free / 1024 / 1024 / 1024),
		DhUsed:    utils.Decimal(curDh.UsedPercent),
		LoadAvg1:  utils.Decimal(loadInfo.Load1),
		LoadAvg5:  utils.Decimal(loadInfo.Load5),
		LoadAvg15: utils.Decimal(loadInfo.Load15),
	}

	return machineUsageInfo
}
