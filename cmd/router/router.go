// SPDX-License-Identifier: Apache-2.0

package router

import (
	"api-routerd/cmd/hostname"
	"api-routerd/cmd/network"
	"api-routerd/cmd/proc"
	"api-routerd/cmd/systemd"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartRouter() {
	router := mux.NewRouter()

	// Register services
	hostname.RegisterRouterHostname(router)
	network.RegisterRouterNetwork(router)
	proc.RegisterRouterProc(router)
	systemd.RegisterRouterSystemd(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
