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

type ProcValue struct {
	Value string `json:"value"`
}

func GetProcNetDev(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetNetDev(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcVersion(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetVersion(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcUserStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetUserStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcTemperatureStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetTemperatureStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcNetStat(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]

	switch r.Method {
	case "GET":
		err := GetNetStat(rw, protocol)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcPidNetStat(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol:= vars["protocol"]
	pid := vars["pid"]

	switch r.Method {
	case "GET":
		err := GetNetStatPid(rw, protocol, pid)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcInterfaceStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetInterfaceStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcProtoCountersStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetProtoCountersStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcGetSwapMemoryStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetSwapMemoryStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcVirtualMemoryStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetVirtualMemoryStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcCPUInfo(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetCPUInfo(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcCPUTimeStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetCPUTimeStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcAvgStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetAvgStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func ConfigureProcSysVM(rw http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)
	vm := ProcVM{Property: vars["path"]}

	switch r.Method {
	case "GET":
		err = vm.GetVM(rw)
		break
	case "PUT":

		v := new(ProcValue)

		err = json.NewDecoder(r.Body).Decode(&v);
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		vm.Value = v.Value
		err = vm.SetVM(rw)
		break
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func ConfigureProcSysNet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := ProcSysNet{Path: vars["path"], Property: vars["conf"] , Link: vars["link"]}

	switch r.Method {
	case "GET":
		err := proc.GetProcSysNet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	case "PUT":

		v := new(ProcValue)
		err := json.NewDecoder(r.Body).Decode(&v);
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		proc.Value = v.Value
		err = proc.SetProcSysNet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcMisc(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetMisc(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcNetArp(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetNetArp(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcModules(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetModules(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func GetProcProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	property := vars["property"]

	switch r.Method {
	case "GET":
		err := GetProcessInfo(rw, pid, property)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter()

	n.HandleFunc("/avgstat", GetProcAvgStat)
	n.HandleFunc("/cpuinfo", GetProcCPUInfo)
	n.HandleFunc("/cputimestat", GetProcCPUTimeStat)
	n.HandleFunc("/interface-stat", GetProcInterfaceStat)
	n.HandleFunc("/misc", GetProcMisc)
	n.HandleFunc("/modules", GetProcModules)
	n.HandleFunc("/net/arp", GetProcNetArp)
	n.HandleFunc("/netdev", GetProcNetDev)
	n.HandleFunc("/netstat/{protocol}", GetProcNetStat)
	n.HandleFunc("/process/{pid}/{method}/", GetProcProcess)
	n.HandleFunc("/proto-counter-stat", GetProcProtoCountersStat)
	n.HandleFunc("/proto-pid-stat/{pid}/{protocol}", GetProcPidNetStat)
	n.HandleFunc("/swap-memory", GetProcGetSwapMemoryStat)
	n.HandleFunc("/sys/net/{path}/{link}/{conf}", ConfigureProcSysNet)
	n.HandleFunc("/sys/vm/{path}", ConfigureProcSysVM)
	n.HandleFunc("/temperaturestat", GetProcTemperatureStat)
	n.HandleFunc("/userstat", GetProcUserStat)
	n.HandleFunc("/version", GetProcVersion)
	n.HandleFunc("/virtual-memory", GetProcVirtualMemoryStat)
}
