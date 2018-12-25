// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/load"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

const ProcMiscPath = "/proc/misc"

func GetVersion(rw http.ResponseWriter) error {
	infostat, r := host.Info()
	if r != nil {
		return r
	}

	b, err := json.Marshal(infostat,)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding: Version")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetNetStat(rw http.ResponseWriter, protocol string) error {
	conn, r := net.Connections(protocol)
	if r != nil {
		return r
	}

	b, err := json.Marshal(conn)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding netstat")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetNetDev(rw http.ResponseWriter) error {
	netdev, r := net.IOCounters(true)
	if r != nil {
		return r
	}

	b, err := json.Marshal(netdev)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding NetDev")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetInterfaceStat(rw http.ResponseWriter) error {
	interfaces, r := net.Interfaces()
	if r != nil {
		return r
	}

	b, err := json.Marshal(interfaces)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding interface stat")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetSwapMemoryStat(rw http.ResponseWriter) error {
	swap, r := mem.SwapMemory()
	if r != nil {
		return r
	}

	b, err := json.Marshal(swap)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding memory stat")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetVirtualMemoryStat(rw http.ResponseWriter) error {
	virt, r := mem.VirtualMemory()
	if r != nil {
		return r
	}

	b, err := json.Marshal(virt)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding VM stat")
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetCPUInfo(rw http.ResponseWriter) error {
	cpus, r := cpu.Info()
	if r != nil {
		return r
	}

	b, err := json.Marshal(cpus)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding CPU Info")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetCPUTimeStat(rw http.ResponseWriter) error {
	cpus, r := cpu.Times(true)
	if r != nil {
		return r
	}

	b, err := json.Marshal(cpus)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding CPU stat")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetAvgStat(rw http.ResponseWriter) error {
	avgstat, r := load.Avg()
	if r != nil {
		return r
	}

	b, err := json.Marshal(avgstat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding avg stat")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}

func GetMisc(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(ProcMiscPath)
	if err != nil {
		log.Fatal("Failed to read: %s", ProcMiscPath)
		return errors.New("Failed to read misc")
	}

	miscMap := make(map[int]string)
	for _, line := range lines {
		fields := strings.Fields(line)

		deviceNum, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		miscMap[deviceNum] = fields[1]
	}

	b, err := json.Marshal(miscMap)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return errors.New("Json encoding")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}
