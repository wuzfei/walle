package common

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/wuzfei/go-helper/unit"
	"runtime"
	"yema.dev/app/utils"
)

type ServerInfo struct {
	User     string `json:"user"`
	Hostname string `json:"hostname"`
	Os       Os     `json:"os"`
	Cpu      Cpu    `json:"cpu"`
	Ram      Ram    `json:"ram"`
	Disk     Disk   `json:"disk"`
}

type Os struct {
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	CPUNum       int    `json:"cpu_num"`
	Compiler     string `json:"compiler"`
	GoVersion    string `json:"go_version"`
	GoroutineNum int    `json:"goroutine_num"`
}

type Cpu struct {
	UsedPercent float64 `json:"used_percent"`
	Cores       int     `json:"cores"`
}

type Ram struct {
	Used        string  `json:"used"`
	Total       string  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

type Disk struct {
	Used        string  `json:"used"`
	Total       string  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

func getServerInfo() (_ *ServerInfo, err error) {
	s := ServerInfo{
		User:     utils.CurrentUser.Username,
		Hostname: utils.CurrentHostname,
	}
	s.Os = initOS()
	if s.Cpu, err = initCPU(); err != nil {
		return
	}
	if s.Ram, err = initRAM(); err != nil {
		return
	}
	if s.Disk, err = initDisk(); err != nil {
		return
	}
	return &s, nil
}

func initOS() (o Os) {
	o.GOOS = runtime.GOOS
	o.GOARCH = runtime.GOARCH
	o.CPUNum = runtime.NumCPU()
	o.Compiler = runtime.Compiler
	o.GoVersion = runtime.Version()
	o.GoroutineNum = runtime.NumGoroutine()
	return o
}

func initCPU() (c Cpu, err error) {
	if cores, err := cpu.Counts(true); err != nil {
		return c, err
	} else {
		c.Cores = cores
	}
	if cpus, err := cpu.Percent(0, false); err != nil {
		return c, err
	} else {
		c.UsedPercent = cpus[0]
	}
	return c, nil
}

func initRAM() (r Ram, err error) {
	if u, err := mem.VirtualMemory(); err != nil {
		return r, err
	} else {
		r.Used = unit.ByteFormat(int64(u.Used), 2)
		r.Total = unit.ByteFormat(int64(u.Total), 2)
		r.UsedPercent = u.UsedPercent
	}
	return r, nil
}

func initDisk() (d Disk, err error) {
	if u, err := disk.Usage("/"); err != nil {
		return d, err
	} else {
		d.Used = unit.ByteFormat(int64(u.Used), 2)
		d.Total = unit.ByteFormat(int64(u.Total), 2)
		d.UsedPercent = u.UsedPercent
	}
	return d, nil
}
