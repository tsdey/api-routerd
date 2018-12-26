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

func GetProc(rw http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)
	path := vars["path"]

	switch r.Method {
	case "GET":
		switch path {
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
			err = GetNetStat(rw, path)
			break
		case "interface-stat":
			err = GetInterfaceStat(rw)
			break
		case "proto-counter-stat":
			err = GetProtoCountersStat(rw)
			break
		case "proto-pid-stat":
			err = GetNetStatPid(rw, path)
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

func ConfigureProcSysVM(rw http.ResponseWriter, r *http.Request) {
	proc := new(ProcVM)

	err := json.NewDecoder(r.Body).Decode(&proc);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		err = proc.GetVM(rw)
		break
	case "PUT":
		err = proc.SetVM(rw)
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
	n.HandleFunc("/{path}", GetProc)
	n.HandleFunc("/misc", GetProcMisc)
	n.HandleFunc("/modules", GetProcModules)
	n.HandleFunc("/net/arp", GetProcNetArp)
	n.HandleFunc("/process/{pid}/{property}", GetProcProcess)
	n.HandleFunc("/sys/vm", ConfigureProcSysVM)
	n.HandleFunc("/sys/net/{path}/{link}/{conf}", ConfigureProcSysNet)
}
