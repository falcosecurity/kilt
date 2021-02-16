package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type ProcessInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Pid int32 `json:"pid"`
	Name string `json:"name"`
	Status string `json:"status"`
	Parent int32 `json:"parent"`
	MemRss uint64 `json:"mem_rss"`
	MemVMS uint64 `json:"mem_vms"`
	CpuUser float64 `json:"cpu_user"`
	CpuSystem float64 `json:"cpu_system"`
}

func main() {
	done := make(chan bool)
	t := time.NewTicker(time.Duration(30) * time.Second)
	defer t.Stop()
	e := json.NewEncoder(os.Stderr)

	go func() {
		for {
			select {
			case <-done:
				return
			case ts := <-t.C:
				ps, _ := process.Processes()
				for _, p := range ps {
					info := getProcessInfo(ts, p)
					e.Encode(info)
				}
			}
		}
	}()

	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)

	<- s
	done <- true
}

func getProcessInfo(ts time.Time, p *process.Process) *ProcessInfo {
	name, err := p.Name()
	if err != nil {
		name = "ERROR: " + err.Error()
	}
	status, err := p.Status()
	if err != nil {
		status = []string{"ERROR:" + err.Error()}
	}

	ppid, _ := p.Ppid()
	memInfo, err := p.MemoryInfo()
	if err!= nil {
		fmt.Printf("error retrieving mem stats for pid %d - %s", p.Pid, err.Error())
		memInfo = &process.MemoryInfoStat{}
	}
	cpuInfo, err := p.Times()
	if err != nil {
		fmt.Printf("error retrieving cpu stats for pid %d - %s", p.Pid, err.Error())
		cpuInfo = &cpu.TimesStat{}
	}

	return &ProcessInfo{
		Timestamp: ts,
		Pid:       p.Pid,
		Name:      name,
		Status:    strings.Join(status, ","),
		Parent:    ppid,
		MemRss:    memInfo.RSS,
		MemVMS:    memInfo.VMS,
		CpuUser:   cpuInfo.User,
		CpuSystem: cpuInfo.System,
	}
}
