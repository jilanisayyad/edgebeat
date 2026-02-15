package controller

import (
	"testing"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/sensors"
)

func TestAvgFloat64(t *testing.T) {
	if got := avgFloat64(nil); got != 0 {
		t.Fatalf("avgFloat64(nil) = %v, want 0", got)
	}
	if got := avgFloat64([]float64{1, 2, 3}); got != 2 {
		t.Fatalf("avgFloat64([1 2 3]) = %v, want 2", got)
	}
}

func TestMapCPUInfo(t *testing.T) {
	in := []cpu.InfoStat{{ModelName: "TestCPU", Cores: 4, Mhz: 3200, CacheSize: 8192}}
	out := mapCPUInfo(in)
	if len(out) != 1 || out[0].ModelName != "TestCPU" || out[0].Cores != 4 {
		t.Fatalf("mapCPUInfo = %+v", out)
	}
}

func TestMapTimes(t *testing.T) {
	in := cpu.TimesStat{User: 1, System: 2, Idle: 3, Nice: 4, Iowait: 5, Irq: 6, Softirq: 7, Steal: 8, Guest: 9, GuestNice: 10}
	out := mapTimes(in)
	if out.User != 1 || out.GuestNice != 10 {
		t.Fatalf("mapTimes = %+v", out)
	}
}

func TestMapTimesSlice(t *testing.T) {
	in := []cpu.TimesStat{{User: 1}, {User: 2}}
	out := mapTimesSlice(in)
	if len(out) != 2 || out[0].User != 1 || out[1].User != 2 {
		t.Fatalf("mapTimesSlice = %+v", out)
	}
}

func TestMapPartitions(t *testing.T) {
	in := []disk.PartitionStat{{Device: "/dev/disk1", Mountpoint: "/", Fstype: "apfs"}}
	out := mapPartitions(in)
	if len(out) != 1 || out[0].Device != "/dev/disk1" || out[0].Mountpoint != "/" {
		t.Fatalf("mapPartitions = %+v", out)
	}
}

func TestMapDiskUsageError(t *testing.T) {
	in := []disk.PartitionStat{{Device: "bad", Mountpoint: "/path/does/not/exist", Fstype: "ext4"}}
	var errs []string
	out := mapDiskUsage(in, &errs)
	if len(out) != 0 {
		t.Fatalf("mapDiskUsage len = %d, want 0", len(out))
	}
	if len(errs) == 0 {
		t.Fatal("expected error for invalid mountpoint")
	}
}

func TestMapDiskIO(t *testing.T) {
	in := map[string]disk.IOCountersStat{"disk0": {ReadBytes: 10, WriteBytes: 20, ReadCount: 1, WriteCount: 2, ReadTime: 3, WriteTime: 4}}
	out := mapDiskIO(in)
	if len(out) != 1 || out[0].Device != "disk0" || out[0].ReadBytes != 10 {
		t.Fatalf("mapDiskIO = %+v", out)
	}
}

func TestMapInterfaces(t *testing.T) {
	in := []net.InterfaceStat{{Name: "eth0", MTU: 1500, HardwareAddr: "00:11:22:33:44:55", Flags: []string{"up"}, Addrs: []net.InterfaceAddr{{Addr: "127.0.0.1"}}}}
	out := mapInterfaces(in)
	if len(out) != 1 || out[0].Name != "eth0" || len(out[0].Addrs) != 1 {
		t.Fatalf("mapInterfaces = %+v", out)
	}
}

func TestMapUsers(t *testing.T) {
	in := []host.UserStat{{User: "tester", Terminal: "pts/0", Host: "localhost", Started: 12345}}
	out := mapUsers(in)
	if len(out) != 1 || out[0].User != "tester" || out[0].StartedUnix != 12345 {
		t.Fatalf("mapUsers = %+v", out)
	}
}

func TestMapTemps(t *testing.T) {
	in := []sensors.TemperatureStat{{SensorKey: "cpu", Temperature: 40, High: 80, Critical: 100}}
	out := mapTemps(in)
	if len(out) != 1 || out[0].SensorKey != "cpu" || out[0].Value != 40 {
		t.Fatalf("mapTemps = %+v", out)
	}
}

func TestCollectSystemInfoNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("collectSystemInfo panicked: %v", r)
		}
	}()

	info := collectSystemInfo()
	if info.Timestamp == "" {
		t.Fatal("collectSystemInfo returned empty timestamp")
	}
	if len(info.Errors) > 0 {
		t.Logf("collectSystemInfo errors: %v", info.Errors)
	}
}
