// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type ProcInfo struct {
	Path string `json:"path"`
	Property string `json:"property"`
}

func GetProc(rw http.ResponseWriter, req *http.Request) {
	proc := new(ProcInfo)

	err := json.NewDecoder(req.Body).Decode(&proc);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "GET":
		switch proc.Path {
		case "netdev":
			err = GetNetDev(rw)
			break
		case "version":
			err = GetVersion(rw)
			break
		case "userstat":
			err = GetUserStat(rw)
			break
		case "temperaturestat":
			err = GetTemperatureStat(rw)
			break
		case "netstat":
			err = GetNetStat(rw, proc.Property)
			break
		case "interface-stat":
			err = GetInterfaceStat(rw)
			break
		case "proto-counter-stat":
			err = GetProtoCountersStat(rw)
			break
		case "proto-pid-stat":
			err = GetNetStatPid(rw, proc.Property)
			break
		case "swap-memory":
			err = GetSwapMemoryStat(rw)
			break
		case "virtual-memory":
			err = GetVirtualMemoryStat(rw)
			break
		case "cpuinfo":
			err = GetCPUInfo(rw)
			break
		case "cputimestat":
			err = GetCPUTimeStat(rw)
			break
		case "avgstat":
			err = GetAvgStat(rw)
			break
		}
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func ConfigureProcSysVM(rw http.ResponseWriter, req *http.Request) {
	proc := new(ProcVM)

	err := json.NewDecoder(req.Body).Decode(&proc);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "GET":
		err = proc.GetVM(rw)
		break
	case "PUT":
		err = proc.SetVM(rw)
		break
	}

	if err != nil {
		http.Error(rw, "Failed to configure VM", http.StatusInternalServerError)
	}
}

func ConfigureProcSysNet(rw http.ResponseWriter, req *http.Request) {
	proc := new(ProcSysNet)

	err := json.NewDecoder(req.Body).Decode(&proc);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "GET":
		err = proc.GetProcSysNet(rw)
		break
	case "PUT":
		err = proc.SetProcSysNet(rw)
		break
	}

	if err != nil {
		http.Error(rw, "Failed to configure /proc/sys/net", http.StatusInternalServerError)
	}
}

func GetProcMisc(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		err := GetMisc(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcNetArp(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		err := GetNetArp(rw)
		if err != nil {
			http.Error(rw, "Failed to get /proc/net/arp", http.StatusInternalServerError)
		}
		break
	}
}

func GetProcModules(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		err := GetModules(rw)
		if err != nil {
			http.Error(rw, "Failed to get /proc/module", http.StatusInternalServerError)
		}
		break
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter()
	n.HandleFunc("/", GetProc)
	n.HandleFunc("/misc", GetProcMisc)
	n.HandleFunc("/modules", GetProcModules)
	n.HandleFunc("/net/arp", GetProcNetArp)
	n.HandleFunc("/sys/vm", ConfigureProcSysVM)
	n.HandleFunc("/sys/net", ConfigureProcSysNet)
}
