/**
 * Author: haoshuaiwei 
 * Date: 2019-05-23 15:36 
 */

package fast_crawl_client

import (
	"fast-work/fast-drive"
	"github.com/hsw409328/gofunc"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"runtime"
	"time"
)

type CpuData struct {
	CpuNumber    int `cpu核数`
	UsedNumber   int `cpu使用数量`
	UsableNumber int `cpu可用数量`
}

func mClient() *host.InfoStat {
	hostInfo, _ := host.Info()
	return hostInfo
}

func mCpu() int {
	return runtime.NumCPU()
}

type MemData struct {
	TotalMem    int
	UsableMem   float64
	UsedPercent float64
}

func mMemory() *MemData {
	memValue, _ := mem.VirtualMemory()
	baseByte := float64(1024 * 1024 * 1024)
	return &MemData{
		TotalMem:    int(float64(memValue.Total) / baseByte),
		UsableMem:   float64(float64(memValue.Available) / baseByte),
		UsedPercent: memValue.UsedPercent,
	}
}

type DiskData struct {
	Total       int
	Used        int
	Usable      int
	UsedPercent float64
}

func mDisk() *DiskData {
	d, _ := disk.Usage("/")
	baseByte := uint64(1024 * 1024 * 1024)
	return &DiskData{
		Total:       int(d.Total / baseByte),
		Used:        int(d.Used / baseByte),
		Usable:      int(d.Free / baseByte),
		UsedPercent: d.UsedPercent,
	}
}

func MonitorClient() {
	go func() {
		for {
			//每10秒汇报一下状态
			hostInfo := mClient()
			memInfo := mMemory()
			diskInfo := mDisk()
			saveResult, err := gofunc.MapOrSliceToJsonString(map[string]interface{}{
				"HostInfo":   hostInfo,
				"MemInfo":    memInfo,
				"DiskInfo":   diskInfo,
				"CreateTime": gofunc.CurrentTime(),
			})
			if err != nil {
				clientLog.Error(err)
			}
			result := fast_drive.RedisDriver.HSet("ClientInfo", hostInfo.Hostname, saveResult)
			log.Println(result)
			time.Sleep(time.Second * 10)
		}
	}()
}
