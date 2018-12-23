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

	r := json.NewDecoder(req.Body).Decode(&hostname);
	if r != nil {
		log.Error("Failed to find decode json: ", r)
		rw.Write([]byte("500: " + r.Error()))
		return
	}

	switch req.Method {
	case "GET":
		switch hostname.Action {
		case "get-hostname":
			GetHostname(rw, hostname)
			break
		}

		break

	case "PUT":
		switch hostname.Action {
		case "set-hostname":
			SetHostname(hostname)
			break
		}
	}
}

func RegisterRouterHostname(router *mux.Router) {
	s := router.PathPrefix("/hostname").Subrouter()
	s.HandleFunc("/", ConfigureHostname)
}
