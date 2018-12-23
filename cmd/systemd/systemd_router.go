// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ConfigureUnit(rw http.ResponseWriter, req *http.Request) {
	unit := new(Unit)

	r := json.NewDecoder(req.Body).Decode(&unit);
	if r != nil {
		log.Error("Failed to find decode json: ", r)
		rw.Write([]byte("500: " + r.Error()))
		return
	}

	switch req.Method {

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

	case "GET":
		GetUnitStatus(rw, unit.Unit)
		break;
	}
}

func ConfigureUnitProperty(rw http.ResponseWriter, req *http.Request) {
	unit := new(Unit)

	r := json.NewDecoder(req.Body).Decode(&unit);
	if r != nil {
		log.Error("Failed to find decode json: ", r)
		rw.Write([]byte("500: " + r.Error()))
		return
	}

	switch req.Method {
	case "GET":
		GetUnitProperty(rw, unit.Unit, unit.Property)
		break
	case "PUT":
		SetUnitProperty(rw, unit.Unit, unit.Property, unit.Value)
		break
	}
}

func RegisterRouterSystemd(router *mux.Router) {
	n := router.PathPrefix("/service").Subrouter()
	n.HandleFunc("/systemd", ConfigureUnit)
	n.HandleFunc("/systemd/property", ConfigureUnitProperty)
}
