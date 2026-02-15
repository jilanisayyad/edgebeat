package controller

import (
	"time"

	"github.com/jilanisayyad/edgebeat/pkg/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/sensors"
)

func collectSystemInfo() utils.SystemInfo {
	errors := make([]string, 0)

	info := utils.SystemInfo{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	}

	if loadAvg, err := load.Avg(); err == nil {
		info.Load = utils.LoadStats{
			Load1:  loadAvg.Load1,
			Load5:  loadAvg.Load5,
			Load15: loadAvg.Load15,
		}
	} else {
		errors = append(errors, "load.Avg: "+err.Error())
	}

	if cpuInfo, err := cpu.Info(); err == nil {
		info.CPU.Info = mapCPUInfo(cpuInfo)
	} else {
		errors = append(errors, "cpu.Info: "+err.Error())
	}

	if perCPUPercent, err := cpu.Percent(0, true); err == nil {
		info.CPU.PerCPUPercent = perCPUPercent
	} else {
		errors = append(errors, "cpu.Percent per-cpu: "+err.Error())
	}

	if totalPercent, err := cpu.Percent(0, false); err == nil && len(totalPercent) > 0 {
		info.CPU.TotalPercent = totalPercent[0]
	} else if len(info.CPU.PerCPUPercent) > 0 {
		info.CPU.TotalPercent = avgFloat64(info.CPU.PerCPUPercent)
	} else if err != nil {
		errors = append(errors, "cpu.Percent total: "+err.Error())
	}

	if totalTimes, err := cpu.Times(false); err == nil && len(totalTimes) > 0 {
		info.CPU.TotalTimes = mapTimes(totalTimes[0])
	} else if err != nil {
		errors = append(errors, "cpu.Times total: "+err.Error())
	}

	if perTimes, err := cpu.Times(true); err == nil {
		info.CPU.PerCPUTimes = mapTimesSlice(perTimes)
	} else {
		errors = append(errors, "cpu.Times per-cpu: "+err.Error())
	}

	if vm, err := mem.VirtualMemory(); err == nil {
		info.Memory.Virtual = utils.VirtualMemory{
			Total:       vm.Total,
			Available:   vm.Available,
			Used:        vm.Used,
			Free:        vm.Free,
			Buffers:     vm.Buffers,
			Cached:      vm.Cached,
			Active:      vm.Active,
			Inactive:    vm.Inactive,
			UsedPercent: vm.UsedPercent,
		}
	} else {
		errors = append(errors, "mem.VirtualMemory: "+err.Error())
	}

	if sm, err := mem.SwapMemory(); err == nil {
		info.Memory.Swap = utils.SwapMemory{
			Total:       sm.Total,
			Used:        sm.Used,
			Free:        sm.Free,
			UsedPercent: sm.UsedPercent,
		}
	} else {
		errors = append(errors, "mem.SwapMemory: "+err.Error())
	}

	if partitions, err := disk.Partitions(false); err == nil {
		info.Disk.Partitions = mapPartitions(partitions)
		info.Disk.Usage = mapDiskUsage(partitions, &errors)
	} else {
		errors = append(errors, "disk.Partitions: "+err.Error())
	}

	if ioStats, err := disk.IOCounters(); err == nil {
		info.Disk.IO = mapDiskIO(ioStats)
	} else {
		errors = append(errors, "disk.IOCounters: "+err.Error())
	}

	if ifaces, err := net.Interfaces(); err == nil {
		info.Network.Interfaces = mapInterfaces(ifaces)
	} else {
		errors = append(errors, "net.Interfaces: "+err.Error())
	}

	if totals, err := net.IOCounters(false); err == nil && len(totals) > 0 {
		info.Network.Totals = utils.NetIO{
			BytesSent:   totals[0].BytesSent,
			BytesRecv:   totals[0].BytesRecv,
			PacketsSent: totals[0].PacketsSent,
			PacketsRecv: totals[0].PacketsRecv,
			Errin:       totals[0].Errin,
			Errout:      totals[0].Errout,
			Dropin:      totals[0].Dropin,
			Dropout:     totals[0].Dropout,
		}
	} else if err != nil {
		errors = append(errors, "net.IOCounters: "+err.Error())
	}

	if hostInfo, err := host.Info(); err == nil {
		info.Host = utils.HostStats{
			Hostname:             hostInfo.Hostname,
			OS:                   hostInfo.OS,
			Platform:             hostInfo.Platform,
			PlatformFamily:       hostInfo.PlatformFamily,
			PlatformVersion:      hostInfo.PlatformVersion,
			KernelVersion:        hostInfo.KernelVersion,
			KernelArch:           hostInfo.KernelArch,
			UptimeSeconds:        hostInfo.Uptime,
			BootTime:             hostInfo.BootTime,
			Procs:                hostInfo.Procs,
			VirtualizationSystem: hostInfo.VirtualizationSystem,
			VirtualizationRole:   hostInfo.VirtualizationRole,
		}
	} else {
		errors = append(errors, "host.Info: "+err.Error())
	}

	if users, err := host.Users(); err == nil {
		info.Host.Users = mapUsers(users)
	} else {
		errors = append(errors, "host.Users: "+err.Error())
	}

	if temps, err := sensors.SensorsTemperatures(); err == nil {
		info.Sensors.Temperatures = mapTemps(temps)
	} else {
		errors = append(errors, "sensors.SensorsTemperatures: "+err.Error())
	}

	if len(errors) > 0 {
		info.Errors = errors
	}

	return info
}

func mapCPUInfo(in []cpu.InfoStat) []utils.CPUInfo {
	items := make([]utils.CPUInfo, 0, len(in))
	for _, v := range in {
		items = append(items, utils.CPUInfo{
			ModelName: v.ModelName,
			Cores:     v.Cores,
			Mhz:       v.Mhz,
			CacheSize: v.CacheSize,
		})
	}
	return items
}

func mapTimes(t cpu.TimesStat) utils.CPUTimes {
	return utils.CPUTimes{
		User:      t.User,
		System:    t.System,
		Idle:      t.Idle,
		Nice:      t.Nice,
		Iowait:    t.Iowait,
		Irq:       t.Irq,
		SoftIrq:   t.Softirq,
		Steal:     t.Steal,
		Guest:     t.Guest,
		GuestNice: t.GuestNice,
	}
}

func mapTimesSlice(in []cpu.TimesStat) []utils.CPUTimes {
	items := make([]utils.CPUTimes, 0, len(in))
	for _, v := range in {
		items = append(items, mapTimes(v))
	}
	return items
}

func mapPartitions(in []disk.PartitionStat) []utils.DiskPartition {
	items := make([]utils.DiskPartition, 0, len(in))
	for _, v := range in {
		items = append(items, utils.DiskPartition{
			Device:     v.Device,
			Mountpoint: v.Mountpoint,
			FSType:     v.Fstype,
		})
	}
	return items
}

func mapDiskUsage(in []disk.PartitionStat, errors *[]string) []utils.DiskUsage {
	items := make([]utils.DiskUsage, 0, len(in))
	for _, v := range in {
		usage, err := disk.Usage(v.Mountpoint)
		if err != nil {
			*errors = append(*errors, "disk.Usage: "+err.Error())
			continue
		}
		items = append(items, utils.DiskUsage{
			Device:      v.Device,
			Mountpoint:  v.Mountpoint,
			FSType:      v.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}
	return items
}

func mapDiskIO(in map[string]disk.IOCountersStat) []utils.DiskIO {
	items := make([]utils.DiskIO, 0, len(in))
	for device, stat := range in {
		items = append(items, utils.DiskIO{
			Device:      device,
			ReadBytes:   stat.ReadBytes,
			WriteBytes:  stat.WriteBytes,
			ReadCount:   stat.ReadCount,
			WriteCount:  stat.WriteCount,
			ReadTimeMS:  stat.ReadTime,
			WriteTimeMS: stat.WriteTime,
		})
	}
	return items
}

func mapInterfaces(in []net.InterfaceStat) []utils.NetInterface {
	items := make([]utils.NetInterface, 0, len(in))
	for _, v := range in {
		addrs := make([]string, 0, len(v.Addrs))
		for _, a := range v.Addrs {
			addrs = append(addrs, a.Addr)
		}
		items = append(items, utils.NetInterface{
			Name:         v.Name,
			MTU:          v.MTU,
			HardwareAddr: v.HardwareAddr,
			Flags:        v.Flags,
			Addrs:        addrs,
		})
	}
	return items
}

func mapUsers(in []host.UserStat) []utils.HostUser {
	items := make([]utils.HostUser, 0, len(in))
	for _, v := range in {
		items = append(items, utils.HostUser{
			User:        v.User,
			Terminal:    v.Terminal,
			Host:        v.Host,
			StartedUnix: int64(v.Started),
		})
	}
	return items
}

func mapTemps(in []sensors.TemperatureStat) []utils.Temperature {
	items := make([]utils.Temperature, 0, len(in))
	for _, v := range in {
		items = append(items, utils.Temperature{
			SensorKey: v.SensorKey,
			Value:     v.Temperature,
			High:      v.High,
			Critical:  v.Critical,
		})
	}
	return items
}

func avgFloat64(in []float64) float64 {
	if len(in) == 0 {
		return 0
	}
	var sum float64
	for _, v := range in {
		sum += v
	}
	return sum / float64(len(in))
}
