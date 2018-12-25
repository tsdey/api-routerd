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

	_= json.NewDecoder(req.Body).Decode(&proc);

	switch req.Method {
	case "GET":
		switch proc.Path {
		case "netdev":
			GetNetDev(rw)
			break
		case "version":
			GetVersion(rw)
			break
		case "netstat":
			GetNetStat(rw, proc.Property)
			break
		case "interface-stat":
			GetInterfaceStat(rw)
			break
		case "swap-memory":
			GetSwapMemoryStat(rw)
			break
		case "virtual-memory":
			GetVirtualMemoryStat(rw)
			break
		case "cpuinfo":
			GetCPUInfo(rw)
			break
		case "cputimestat":
			GetCPUTimeStat(rw)
			break
		case "avgstat":
			GetAvgStat(rw)
			break
		}
	}
}

func ConfigureProcSysVM(rw http.ResponseWriter, req *http.Request) {
	proc := new(ProcVM)

	_= json.NewDecoder(req.Body).Decode(&proc);

	switch req.Method {
	case "GET":
		proc.GetVM(rw)
		break
	case "PUT":
		proc.SetVM(rw)
		break
	}
}

func ConfigureProcSysNet(rw http.ResponseWriter, req *http.Request) {
	proc := new(ProcSysNet)

	_= json.NewDecoder(req.Body).Decode(&proc);

	switch req.Method {
	case "GET":
		proc.GetProcSysNet(rw)
		break
	case "PUT":
		proc.SetProcSysNet(rw)
		break
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter()
	n.HandleFunc("/", GetProc)
	n.HandleFunc("/sys/vm", ConfigureProcSysVM)
	n.HandleFunc("/sys/net", ConfigureProcSysNet)
}
