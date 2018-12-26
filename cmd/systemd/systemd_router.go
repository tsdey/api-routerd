// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func RouterConfigureUnit(rw http.ResponseWriter, r *http.Request) {
	unit := new(Unit)

	err := json.NewDecoder(r.Body).Decode(&unit);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		switch unit.Action {
		case "start":
			StartUnit(unit.Unit)
			break
		case "stop":
			StopUnit(unit.Unit)
			break
		case "restart":
			RestartUnit(unit.Unit)
			break
		case "reload":
			ReloadUnit(unit.Unit)
			break
		}
		break
	}
}

func RouterGetUnitStatus(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unit := vars["unit"]

	switch r.Method {
	case "GET":
		GetUnitStatus(rw, unit)
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
	n.HandleFunc("/systemd", RouterConfigureUnit)
	n.HandleFunc("/systemd/{unit}/status", RouterGetUnitStatus)
	n.HandleFunc("/systemd/{unit}/get/{property}", RouterGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/set/{property}", RouterConfigureUnitProperty)
}
