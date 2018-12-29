// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func RouterConfigureUnit(rw http.ResponseWriter, r *http.Request) {
	var err error

	switch r.Method {
	case "GET":
		err = SystemState(rw)
		break
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
		}
		break
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
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
	n.HandleFunc("/systemd", RouterConfigureUnit)
	n.HandleFunc("/systemd/{unit}/status", RouterGetUnitStatus)
	n.HandleFunc("/systemd/{unit}/get/{property}", RouterGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/set/{property}", RouterConfigureUnitProperty)
}
