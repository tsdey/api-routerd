// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ConfigureHostname(rw http.ResponseWriter, req *http.Request) {
	hostname := new(Hostname)

	err := json.NewDecoder(req.Body).Decode(&hostname);
	if err != nil {
		log.Error("Failed to decode json: ", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Method {
	case "GET":
		err = hostname.GetHostname(rw)
		break

	case "PUT":
		err = hostname.SetHostname()
		break
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func RegisterRouterHostname(router *mux.Router) {
	s := router.PathPrefix("/hostname").Subrouter().StrictSlash(false)
	s.HandleFunc("/", ConfigureHostname)
}
