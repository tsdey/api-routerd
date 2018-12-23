// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/load"
	"net/http"
)

func GetVersion(rw http.ResponseWriter) {
	infostat, r := host.Info()
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(infostat)
}

func GetNetStat(rw http.ResponseWriter, protocol string) {
	conn, r := net.Connections(protocol)
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(conn)
}

func GetNetDev(rw http.ResponseWriter) {
	netdev, r := net.IOCounters(true)
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(netdev)
}

func GetInterfaceStat(rw http.ResponseWriter) {
	interfaces, r := net.Interfaces()
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(interfaces)
}

func GetSwapMemoryStat(rw http.ResponseWriter) {
	swap, r := mem.SwapMemory()
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(swap)
}

func GetVirtualMemoryStat(rw http.ResponseWriter) {
	virt, r := mem.VirtualMemory()
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(virt)
}

func GetCPUInfo(rw http.ResponseWriter) {
	cpus, r := cpu.Info()
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(cpus)
}

func GetCPUTimeStat(rw http.ResponseWriter) {
	cpus, r := cpu.Times(true)
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(cpus)
}

func GetAvgStat(rw http.ResponseWriter) {
	avgstat, r := load.Avg()
	if r != nil {
		return
	}

	json.NewEncoder(rw).Encode(avgstat)
}
