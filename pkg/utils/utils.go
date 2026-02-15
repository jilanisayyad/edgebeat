package utils

type SystemInfo struct {
	Timestamp string       `json:"timestamp"`
	CPU       CPUStats     `json:"cpu"`
	Load      LoadStats    `json:"load"`
	Memory    MemoryStats  `json:"memory"`
	Disk      DiskStats    `json:"disk"`
	Network   NetworkStats `json:"network"`
	Host      HostStats    `json:"host"`
	Sensors   SensorsStats `json:"sensors"`
	Errors    []string     `json:"errors,omitempty"`
}

type CPUStats struct {
	Info          []CPUInfo  `json:"info"`
	TotalTimes    CPUTimes   `json:"total_times"`
	PerCPUTimes   []CPUTimes `json:"per_cpu_times"`
	TotalPercent  float64    `json:"total_percent"`
	PerCPUPercent []float64  `json:"per_cpu_percent"`
}

type CPUInfo struct {
	ModelName string  `json:"model_name"`
	Cores     int32   `json:"cores"`
	Mhz       float64 `json:"mhz"`
	CacheSize int32   `json:"cache_size"`
}

type CPUTimes struct {
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	SoftIrq   float64 `json:"soft_irq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guest_nice"`
}

type LoadStats struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type MemoryStats struct {
	Virtual VirtualMemory `json:"virtual"`
	Swap    SwapMemory    `json:"swap"`
}

type VirtualMemory struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
	Active      uint64  `json:"active"`
	Inactive    uint64  `json:"inactive"`
	UsedPercent float64 `json:"used_percent"`
}

type SwapMemory struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskStats struct {
	Partitions []DiskPartition `json:"partitions"`
	Usage      []DiskUsage     `json:"usage"`
	IO         []DiskIO        `json:"io"`
}

type DiskPartition struct {
	Device     string `json:"device"`
	Mountpoint string `json:"mountpoint"`
	FSType     string `json:"fs_type"`
}

type DiskUsage struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	FSType      string  `json:"fs_type"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskIO struct {
	Device      string `json:"device"`
	ReadBytes   uint64 `json:"read_bytes"`
	WriteBytes  uint64 `json:"write_bytes"`
	ReadCount   uint64 `json:"read_count"`
	WriteCount  uint64 `json:"write_count"`
	ReadTimeMS  uint64 `json:"read_time_ms"`
	WriteTimeMS uint64 `json:"write_time_ms"`
}

type NetworkStats struct {
	Interfaces []NetInterface `json:"interfaces"`
	Totals     NetIO          `json:"totals"`
}

type NetInterface struct {
	Name         string   `json:"name"`
	MTU          int      `json:"mtu"`
	HardwareAddr string   `json:"hardware_addr"`
	Flags        []string `json:"flags"`
	Addrs        []string `json:"addrs"`
}

type NetIO struct {
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	Errin       uint64 `json:"err_in"`
	Errout      uint64 `json:"err_out"`
	Dropin      uint64 `json:"drop_in"`
	Dropout     uint64 `json:"drop_out"`
}

type HostStats struct {
	Hostname             string     `json:"hostname"`
	OS                   string     `json:"os"`
	Platform             string     `json:"platform"`
	PlatformFamily       string     `json:"platform_family"`
	PlatformVersion      string     `json:"platform_version"`
	KernelVersion        string     `json:"kernel_version"`
	KernelArch           string     `json:"kernel_arch"`
	UptimeSeconds        uint64     `json:"uptime_seconds"`
	BootTime             uint64     `json:"boot_time"`
	Procs                uint64     `json:"procs"`
	VirtualizationSystem string     `json:"virtualization_system"`
	VirtualizationRole   string     `json:"virtualization_role"`
	Users                []HostUser `json:"users"`
}

type HostUser struct {
	User        string `json:"user"`
	Terminal    string `json:"terminal"`
	Host        string `json:"host"`
	StartedUnix int64  `json:"started_unix"`
}

type SensorsStats struct {
	Temperatures []Temperature `json:"temperatures"`
	Fans         []Fan         `json:"fans"`
}

type Temperature struct {
	SensorKey string  `json:"sensor_key"`
	Value     float64 `json:"value"`
	High      float64 `json:"high"`
	Critical  float64 `json:"critical"`
}

type Fan struct {
	SensorKey string  `json:"sensor_key"`
	Value     float64 `json:"value"`
}
