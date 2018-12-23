// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"net/http"
	"path"
)

const VMPath = "/proc/sys/vm"

type ProcVM struct {
	Path string `json:"path"`
	Property string `json:"property"`
	Value string `json:"value"`
}

func GetVM(rw http.ResponseWriter, property string) (error) {
	line, r := share.ReadOneLineFile(path.Join(VMPath, property))
	if r != nil {
		return r
	}

	vmproperty := ProcVM {Property: property, Value: line}
	json.NewEncoder(rw).Encode(vmproperty)

	return nil
}

func SetVM(rw http.ResponseWriter, property string, value string) (error) {
	r := share.WriteOneLineFile(path.Join(VMPath, property), value)
	if r != nil {
		return r
	}

	line, r := share.ReadOneLineFile(path.Join(VMPath, property))
	if r != nil {
		return r
	}

	vmproperty := ProcVM {Path: "vm", Property: property, Value: line}
	json.NewEncoder(rw).Encode(vmproperty)

	return nil
}
