// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"errors"
	"net/http"
	"path"
)

const (
	SysNetPath     = "/proc/sys/net"
	SysNetPathCore = "core"
	SysNetPathIPv4 = "ipv4"
	SysNetPathIPv6 = "ipv6"
)

type ProcSysNet struct {
	Path string     `json:"path"`
	Property string `json:"property"`
	Value string    `json:"value"`
	Link string     `json:"link"`
}

func (req *ProcSysNet) GetProcSysNetPath() (string, error) {
	var procPath string

	switch req.Path {
	case SysNetPathCore:
		procPath = path.Join(path.Join(SysNetPath, SysNetPathCore), req.Property)
		break
	case SysNetPathIPv4:

		if (req.Link != "") {
			procPath = path.Join(path.Join(path.Join(path.Join(SysNetPath, SysNetPathIPv4), "conf"), req.Link), req.Property)
		} else {
			procPath = path.Join(path.Join(SysNetPath, SysNetPathIPv4), req.Property)
		}
		break
	case SysNetPathIPv6:

		if (req.Link != "") {
			procPath = path.Join(path.Join(path.Join(path.Join(SysNetPath, SysNetPathIPv6), "conf"), req.Link), req.Property)
		} else {
			procPath = path.Join(path.Join(SysNetPath, SysNetPathIPv6), req.Property)
		}
		break
	default:
		return "", errors.New("Path not found")
	}

	return procPath, nil
}

func (req *ProcSysNet) GetProcSysNet(rw http.ResponseWriter) (error) {
	path, err := req.GetProcSysNetPath()
	if err != nil {
		return err
	}

	line, err := share.ReadOneLineFile(path)
	if err != nil {
		return err
	}

	property := ProcSysNet {Path: req.Path, Property: req.Property, Value: line, Link: req.Link}
	json.NewEncoder(rw).Encode(property)

	return nil
}

func (req *ProcSysNet) SetProcSysNet(rw http.ResponseWriter) (error) {
	path, err := req.GetProcSysNetPath()
	if err != nil {
		return err
	}

	err = share.WriteOneLineFile(path, req.Value)
	if err != nil {
		return err
	}

	return nil
}
