// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func RouterGetHostname(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	property := vars["property"]

	switch r.Method {
	case "GET":
		err := GetHostname(rw, property)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		break
	}
}

func RouterSetHostname(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	property := vars["property"]

	hostname := new(Hostname)

	err := json.NewDecoder(r.Body).Decode(&hostname);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "PUT":
		hostname.Property = property

		err := hostname.SetHostname()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RegisterRouterHostname(router *mux.Router) {
	s := router.PathPrefix("/hostname").Subrouter().StrictSlash(false)
	s.HandleFunc("/get/{property}", RouterGetHostname)
	s.HandleFunc("/set/{property}", RouterSetHostname)
}
