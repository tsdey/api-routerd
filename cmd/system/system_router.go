// SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RouterConfigureJournalConf(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetJournalConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break

	case "POST":
		err := UpdateJournalConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func RegisterRouterSystem(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

	// conf
	n.HandleFunc("/journal/conf", RouterConfigureJournalConf)
	n.HandleFunc("/journal/conf/update", RouterConfigureJournalConf)

}
