// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func RouterGetSystemdState(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := SystemdState(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetSystemdVersion(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := SystemdVersion(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetSystemdFeatures(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := SystemdFeatures(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetSystemdVirtualization(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := SystemdVirtualization(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterGetSystemdArchitecture(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := SystemdArchitecture(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RouterConfigureUnit(rw http.ResponseWriter, r *http.Request) {
	var err error

	switch r.Method {
	case "POST":
		unit := new(Unit)

		err = json.NewDecoder(r.Body).Decode(&unit);
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		switch unit.Action {
		case "start":
			err = unit.StartUnit()
			break
		case "stop":
			err = unit.StopUnit()
			break
		case "restart":
			err = unit.RestartUnit()
			break
		case "reload":
			err = unit.ReloadUnit()
			break
		case "kill":
			err = unit.KillUnit()
			break
		}
		break
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func RouterGetAllUnits(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := ListUnits(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		break;
	}

}

func RouterGetUnitStatus(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unit := vars["unit"]

	u := Unit{Unit: unit}

	switch r.Method {
	case "GET":
		err := u.GetUnitStatus(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		break;
	}
}

func RouterGetUnitProperty(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unit := vars["unit"]
	property := vars["property"]

	u := new(Unit)
	u.Unit = unit
	u.Property = property

	switch r.Method {
	case "GET":
		u.GetUnitProperty(rw)
		break
	}
}

func RouterConfigureUnitProperty(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unit := vars["unit"]
	property := vars["property"]

	u := new(Unit)
	err := json.NewDecoder(r.Body).Decode(&u);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	u.Unit = unit
	u.Property = property

	switch r.Method {
	case "PUT":
		u.SetUnitProperty(rw)
		break
	}
}

func RegisterRouterSystemd(router *mux.Router) {
	n := router.PathPrefix("/service").Subrouter()
	n.HandleFunc("/systemd/state", RouterGetSystemdState)
	n.HandleFunc("/systemd/version", RouterGetSystemdVersion)
	n.HandleFunc("/systemd/features", RouterGetSystemdFeatures)
	n.HandleFunc("/systemd/virtualization", RouterGetSystemdVirtualization)
	n.HandleFunc("/systemd/architecture", RouterGetSystemdArchitecture)
	n.HandleFunc("/systemd/units", RouterGetAllUnits)

	n.HandleFunc("/systemd", RouterConfigureUnit)
	n.HandleFunc("/systemd/{unit}/status", RouterGetUnitStatus)
	n.HandleFunc("/systemd/{unit}/get", RouterGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/get/{property}", RouterGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/set/{property}", RouterConfigureUnitProperty)
}
